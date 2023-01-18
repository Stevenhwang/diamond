<template>
  <div class="dashboard-container">
    <div class="dashboard-text">name: {{ name }}</div>
    <div class="dashboard-text">menus: {{ menus }}</div>
    <el-alert title="1.基于 websocket 的 web terminal，在 server 界面点击连接即可"
              effect="dark"
              :closable="false"
              type="success">
    </el-alert>
    <el-divider></el-divider>
    <el-alert title="2.基于 golang 的高性能 ssh 服务器，支持任意 ssh 客户端连接，支持 navicate ssh 隧道"
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
    <el-alert title="4.方便高效的任务平台，基于 linux command， 内置 ansible 和常用 playbooks"
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
