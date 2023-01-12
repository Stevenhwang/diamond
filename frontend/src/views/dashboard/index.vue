<template>
  <div class="dashboard-container">
    <div class="dashboard-text">name: {{ name }}</div>
    <div class="dashboard-text">menus: {{ menus }}</div>
    <el-alert title="1.支持 web terminal，可以在 server 界面点击连接，就可以连接服务器"
              effect="dark"
              :closable="false"
              type="success">
    </el-alert>
    <el-divider></el-divider>
    <el-alert title="2.支持 ssh 客户端连接(默认端口2222)，比如使用xshell: ssh test1@127.0.0.1 2222"
              effect="dark"
              :closable="false"
              type="success">
    </el-alert>
    <el-divider></el-divider>
    <el-alert title="3.ssh 客户端连接同时支持密码和 public key 认证，安全起见，请联系管理员绑定 key"
              effect="dark"
              :closable="false"
              type="success">
    </el-alert>
    <el-divider></el-divider>
    <el-alert title="4.支持 navicate http 隧道代理，使用时请使用登录时使用的账号和密码"
              effect="dark"
              :closable="false"
              type="success">
    </el-alert>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'Dashboard',
  computed: {
    ...mapGetters([
      'name',
      'menus'
    ])
  }
}
</script>

<style lang="scss" scoped>
.dashboard {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
