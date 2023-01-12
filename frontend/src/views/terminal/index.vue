<template>
  <div class="container">
    <div id="terminal"
         ref="terminal" />
  </div>
</template>

<script>
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { AttachAddon } from 'xterm-addon-attach'
import 'xterm/css/xterm.css'
import { getToken } from '@/utils/auth'

export default {
  data() {
    return {
      scoket: '',
      term: null
    }
  },
  mounted() {
    // 实例化终端并设置参数
    this.term = new Terminal({
      cursorBlink: true
    })
    // 加载自适应组件
    this.fitAddon = new FitAddon()
    this.term.loadAddon(this.fitAddon)

    // 加载weblink组件
    // this.term.loadAddon(new WebLinksAddon())

    // 在绑定的组件上初始化窗口
    this.term.open(this.$refs.terminal)

    // 窗口初始化后,按照浏览器窗口自适应大小
    this.fitAddon.fit()

    // 聚焦
    this.term.focus()

    // 创建ws实例
    // 这里还把窗口的column和row传入后端,使其能自动针对前端窗口边框改为输出
    const host = window.location.host
    const protocal = window.location.protocol
    const head = protocal === "http:" ? "ws:" : "wss:"
    const wsURL = `${head}//${host}/api/terminal`
    const token = getToken()
    const url = `${wsURL}?token=${token}&id=${this.$route.query.id}&cols=${this.term.cols}&rows=${this.term.rows}`
    this.socket = new WebSocket(url)

    // xterm的socket组件与websocket实例结合
    const attachAddon = new AttachAddon(this.socket)
    this.term.loadAddon(attachAddon)

    // 监听resize,当窗口拖动的时候,监听事件,实时改变xterm窗口
    window.addEventListener('resize', this.debounce(this.resizeScreen, 1000), false)
  },
  methods: {
    // 节流,避免拖动时候频繁向后端请求更新
    debounce(fn, wait) {
      let timeout = null
      return function() {
        if (timeout !== null) clearTimeout(timeout)
        timeout = setTimeout(fn, wait)
      }
    },
    // 页面重新resize的时候,需要重新告诉后端cols和rows
    resizeScreen() {
      // 对xterm的窗口重新fit,获取新的cols和rows
      this.fitAddon.fit()
      this.socket.send(JSON.stringify([this.term.cols, this.term.rows]) + '\n')
    }
  }
}
</script>
<style scoped>
.container #terminal {
  height: 99.99%;
  width: 99.99%;
  position: absolute;
  top: 0px;
  bottom: 0px;
  z-index: 0;
  overflow: hidden;
}
</style>
