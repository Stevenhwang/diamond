package socks5

import (
	"context"
	"diamond/cache"
	"diamond/misc"
	"diamond/models"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/Stevenhwang/gommon/tools"
)

const (
	userpassAuthRequired byte = 2
	noAcceptableAuth     byte = 255
	authSuccess          byte = 0
	authVersion          byte = 1

	// socks5Version is the byte that represents the SOCKS version
	// in requests.
	socks5Version byte = 5
)

// commandType are the bytes sent in SOCKS5 packets
// that represent the kind of connection the client needs.
type commandType byte

// The set of valid SOCKS5 commands as described in RFC 1928.
const (
	connect      commandType = 1
	bind         commandType = 2
	udpAssociate commandType = 3
)

// addrType are the bytes sent in SOCKS5 packets
// that represent particular address types.
type addrType byte

// The set of valid SOCKS5 address types as defined in RFC 1928.
const (
	ipv4       addrType = 1
	domainName addrType = 3
	ipv6       addrType = 4
)

// replyCode are the bytes sent in SOCKS5 packets
// that represent replies from the server to a client
// request.
type replyCode byte

// The set of valid SOCKS5 reply types as per the RFC 1928.
const (
	success              replyCode = 0
	generalFailure       replyCode = 1
	connectionNotAllowed replyCode = 2
	networkUnreachable   replyCode = 3
	hostUnreachable      replyCode = 4
	connectionRefused    replyCode = 5
	ttlExpired           replyCode = 6
	commandNotSupported  replyCode = 7
	addrTypeNotSupported replyCode = 8
)

// Server is a SOCKS5 proxy server.
type Server struct {
}

func (s *Server) dial(ctx context.Context, network, addr string) (net.Conn, error) {
	dialer := &net.Dialer{}
	dial := dialer.DialContext
	return dial(ctx, network, addr)
}

// Serve accepts and handles incoming connections on the given listener.
func (s *Server) Serve(l net.Listener) error {
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}
		go func() {
			defer c.Close()
			conn := &Conn{clientConn: c, srv: s}
			err := conn.Run()
			if err != nil {
				misc.Logger.Error().Err(fmt.Errorf("client connection failed: %v", err)).Str("from", "socks5").Msg("")
			}
		}()
	}
}

// Conn is a SOCKS5 connection for client to reach
// server.
type Conn struct {
	// The struct is filled by each of the internal
	// methods in turn as the transaction progresses.

	srv        *Server
	clientConn net.Conn
	request    *request
}

// 处理握手阶段
func handshake(r io.Reader) error {
	var hdr [2]byte
	_, err := io.ReadFull(r, hdr[:])
	if err != nil {
		return fmt.Errorf("could not read packet header")
	}
	if hdr[0] != socks5Version {
		return fmt.Errorf("incompatible SOCKS version")
	}
	count := int(hdr[1])
	methods := make([]byte, count)
	_, err = io.ReadFull(r, methods)
	if err != nil {
		return fmt.Errorf("could not read methods")
	}
	for _, m := range methods {
		if m == userpassAuthRequired {
			return nil
		}
	}
	return fmt.Errorf("no acceptable auth methods")
}

// 处理 user pass 认证
func auth(remoteip string, r io.Reader) error {
	var hdr [2]byte
	_, err := io.ReadFull(r, hdr[:])
	if err != nil {
		return fmt.Errorf("could not read packet header")
	}
	if hdr[0] != authVersion {
		return fmt.Errorf("incompatible SOCKS version")
	}
	userLen := int(hdr[1])
	username := make([]byte, userLen)
	_, err = io.ReadFull(r, username)
	if err != nil {
		return fmt.Errorf("could not read username")
	}
	misc.Logger.Info().Str("from", "socks5").Msg(fmt.Sprintf("get user: %s", string(username)))
	// 查找数据库有没有此用户
	user := models.User{}
	if res := models.DB.Where("username = ?", string(username)).First(&user); res.Error != nil {
		cache.Ban(remoteip)
		return fmt.Errorf("username not valid")
	}
	// 如果用户被禁用
	if !user.IsActive {
		return fmt.Errorf("user is not active")
	}
	var tail [1]byte
	_, err = io.ReadFull(r, tail[:])
	if err != nil {
		return fmt.Errorf("could not read pass len")
	}
	passLen := int(tail[0])
	pass := make([]byte, passLen)
	_, err = io.ReadFull(r, pass)
	if err != nil {
		return fmt.Errorf("could not read pass")
	}
	// 验证密码
	if tools.CheckPassword(user.Password, string(pass)) {
		return nil
	} else {
		cache.Ban(remoteip)
		return fmt.Errorf("password not valid")
	}
}

// Run starts the new connection.
func (c *Conn) Run() error {
	err := handshake(c.clientConn)
	if err != nil {
		c.clientConn.Write([]byte{socks5Version, noAcceptableAuth})
		return err
	}
	c.clientConn.Write([]byte{socks5Version, userpassAuthRequired})
	remoteip := c.clientConn.RemoteAddr().(*net.TCPAddr).IP.String()
	if err := auth(remoteip, c.clientConn); err != nil {
		c.clientConn.Write([]byte{socks5Version, noAcceptableAuth})
		return err
	}
	c.clientConn.Write([]byte{socks5Version, authSuccess})
	return c.handleRequest()
}

func (c *Conn) handleRequest() error {
	req, err := parseClientRequest(c.clientConn)
	if err != nil {
		res := &response{reply: generalFailure}
		buf, _ := res.marshal()
		c.clientConn.Write(buf)
		return err
	}
	if req.command != connect {
		res := &response{reply: commandNotSupported}
		buf, _ := res.marshal()
		c.clientConn.Write(buf)
		return fmt.Errorf("unsupported command %v", req.command)
	}
	c.request = req

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv, err := c.srv.dial(
		ctx,
		"tcp",
		net.JoinHostPort(c.request.destination, strconv.Itoa(int(c.request.port))),
	)
	if err != nil {
		res := &response{reply: generalFailure}
		buf, _ := res.marshal()
		c.clientConn.Write(buf)
		return err
	}
	defer srv.Close()
	serverAddr, serverPortStr, err := net.SplitHostPort(srv.LocalAddr().String())
	if err != nil {
		return err
	}
	serverPort, _ := strconv.Atoi(serverPortStr)

	var bindAddrType addrType
	if ip := net.ParseIP(serverAddr); ip != nil {
		if ip.To4() != nil {
			bindAddrType = ipv4
		} else {
			bindAddrType = ipv6
		}
	} else {
		bindAddrType = domainName
	}
	res := &response{
		reply:        success,
		bindAddrType: bindAddrType,
		bindAddr:     serverAddr,
		bindPort:     uint16(serverPort),
	}
	buf, err := res.marshal()
	if err != nil {
		res = &response{reply: generalFailure}
		buf, _ = res.marshal()
	}
	c.clientConn.Write(buf)

	errc := make(chan error, 2)
	go func() {
		_, err := io.Copy(c.clientConn, srv)
		if err != nil {
			err = fmt.Errorf("from backend to client: %w", err)
		}
		errc <- err
	}()
	go func() {
		_, err := io.Copy(srv, c.clientConn)
		if err != nil {
			err = fmt.Errorf("from client to backend: %w", err)
		}
		errc <- err
	}()
	return <-errc
}

// request represents data contained within a SOCKS5
// connection request packet.
type request struct {
	command      commandType
	destination  string
	port         uint16
	destAddrType addrType
}

// parseClientRequest converts raw packet bytes into a
// SOCKS5Request struct.
func parseClientRequest(r io.Reader) (*request, error) {
	var hdr [4]byte
	_, err := io.ReadFull(r, hdr[:])
	if err != nil {
		return nil, fmt.Errorf("could not read packet header")
	}
	cmd := hdr[1]
	destAddrType := addrType(hdr[3])

	var destination string
	var port uint16

	if destAddrType == ipv4 {
		var ip [4]byte
		_, err = io.ReadFull(r, ip[:])
		if err != nil {
			return nil, fmt.Errorf("could not read IPv4 address")
		}
		destination = net.IP(ip[:]).String()
	} else if destAddrType == domainName {
		var dstSizeByte [1]byte
		_, err = io.ReadFull(r, dstSizeByte[:])
		if err != nil {
			return nil, fmt.Errorf("could not read domain name size")
		}
		dstSize := int(dstSizeByte[0])
		domainName := make([]byte, dstSize)
		_, err = io.ReadFull(r, domainName)
		if err != nil {
			return nil, fmt.Errorf("could not read domain name")
		}
		destination = string(domainName)
	} else if destAddrType == ipv6 {
		var ip [16]byte
		_, err = io.ReadFull(r, ip[:])
		if err != nil {
			return nil, fmt.Errorf("could not read IPv6 address")
		}
		destination = net.IP(ip[:]).String()
	} else {
		return nil, fmt.Errorf("unsupported address type")
	}
	var portBytes [2]byte
	_, err = io.ReadFull(r, portBytes[:])
	if err != nil {
		return nil, fmt.Errorf("could not read port")
	}
	port = binary.BigEndian.Uint16(portBytes[:])

	return &request{
		command:      commandType(cmd),
		destination:  destination,
		port:         port,
		destAddrType: destAddrType,
	}, nil
}

// response contains the contents of
// a response packet sent from the proxy
// to the client.
type response struct {
	reply        replyCode
	bindAddrType addrType
	bindAddr     string
	bindPort     uint16
}

// marshal converts a SOCKS5Response struct into
// a packet. If res.reply == Success, it may throw an error on
// receiving an invalid bind address. Otherwise, it will not throw.
func (res *response) marshal() ([]byte, error) {
	pkt := make([]byte, 4)
	pkt[0] = socks5Version
	pkt[1] = byte(res.reply)
	pkt[2] = 0 // null reserved byte
	pkt[3] = byte(res.bindAddrType)

	if res.reply != success {
		return pkt, nil
	}

	var addr []byte
	switch res.bindAddrType {
	case ipv4:
		addr = net.ParseIP(res.bindAddr).To4()
		if addr == nil {
			return nil, fmt.Errorf("invalid IPv4 address for binding")
		}
	case domainName:
		if len(res.bindAddr) > 255 {
			return nil, fmt.Errorf("invalid domain name for binding")
		}
		addr = make([]byte, 0, len(res.bindAddr)+1)
		addr = append(addr, byte(len(res.bindAddr)))
		addr = append(addr, []byte(res.bindAddr)...)
	case ipv6:
		addr = net.ParseIP(res.bindAddr).To16()
		if addr == nil {
			return nil, fmt.Errorf("invalid IPv6 address for binding")
		}
	default:
		return nil, fmt.Errorf("unsupported address type")
	}

	pkt = append(pkt, addr...)
	pkt = binary.BigEndian.AppendUint16(pkt, uint16(res.bindPort))

	return pkt, nil
}

func Start(addr string) {
	s := Server{}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Fatal(s.Serve(l))
}
