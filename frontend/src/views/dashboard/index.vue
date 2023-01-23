<template>
  <div class="dashboard-container">
    <div class="dashboard-text">Hello, {{ name }}</div>
    <!-- <div class="dashboard-text">menus: {{ menus }}</div> -->
    <el-alert title="1.基于 websocket 的 web terminal，在 server 菜单页点击相应服务器的连接按钮即可使用"
              :closable="false"
              type="info">
    </el-alert>
    <el-divider></el-divider>
    <el-alert title="2.golang 开发的高性能 ssh 服务器，支持任意 ssh 客户端连接，支持 navicate ssh 隧道(使用平台账密或key)"
              :closable="false"
              type="info">
    </el-alert>
    <el-divider></el-divider>
    <el-alert title="3.ssh 客户端连接(如 Xshell 或 SecureCRT 等)可同时支持密码和 key，如需要请联系管理员绑定个人 key"
              :closable="false"
              type="info">
    </el-alert>
    <el-divider></el-divider>
    <el-alert title="4.方便高效的任务平台，基于 linux command， 内置 ansible 和常用 playbooks，可处理常见的部署和批处理任务"
              :closable="false"
              type="info">
    </el-alert>
    <el-divider></el-divider>
    <el-alert title="5.在任务平台基础上，建造了强大的定时任务模块，支持秒级的crontab表达式"
              :closable="false"
              type="info">
    </el-alert>
    <el-divider></el-divider>
    <el-alert title="6.golang 开发的高性能 socks5 代理服务器(使用平台账密连接)，便于跟服务器内网通讯，比如使用 ftp 客户端等"
              :closable="false"
              type="info">
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
      // 'menus'
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
