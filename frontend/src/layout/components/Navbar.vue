<template>
  <div class="navbar">
    <hamburger :is-active="sidebar.opened"
               class="hamburger-container"
               @toggleClick="toggleSideBar" />

    <breadcrumb class="breadcrumb-container" />

    <!-- 修改密码弹框 -->
    <el-dialog title="修改密码"
               :visible.sync="dialogFormVisible"
               width="35%">
      <el-form ref="form"
               :model="form">
        <el-form-item label="原始密码:"
                      prop="old_pw">
          <el-input v-model="form.old_pw"
                    type="password"
                    placeholder="请输入原始密码"
                    autocomplete="off" />
        </el-form-item>
        <el-form-item label="新密码:"
                      prop="new_pw1">
          <el-input v-model="form.new_pw1"
                    type="password"
                    placeholder="请输入新密码"
                    autocomplete="off" />
        </el-form-item>
        <el-form-item label="确认密码:"
                      prop="new_pw2">
          <el-input v-model="form.new_pw2"
                    type="password"
                    placeholder="请再次输入密码"
                    autocomplete="off" />
        </el-form-item>
      </el-form>
      <div slot="footer"
           class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button type="primary"
                   @click="modify">确 定</el-button>
      </div>
    </el-dialog>
    <!-- 修改密码框 -->

    <div class="right-menu">
      <el-dropdown class="avatar-container"
                   trigger="click">
        <div class="avatar-wrapper">
          <img :src="avatar+'?imageView2/1/w/80/h/80'"
               class="user-avatar">
          <i class="el-icon-caret-bottom" />
        </div>
        <el-dropdown-menu slot="dropdown"
                          class="user-dropdown">
          <router-link to="/">
            <el-dropdown-item>
              Home
            </el-dropdown-item>
          </router-link>
          <el-dropdown-item divided
                            @click.native="changePass">
            <span style="display:block;">Change Pass</span>
          </el-dropdown-item>
          <el-dropdown-item divided
                            @click.native="logout">
            <span style="display:block;">Log Out</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import Breadcrumb from '@/components/Breadcrumb'
import Hamburger from '@/components/Hamburger'
import { resetPass } from '@/api/user'

export default {
  components: {
    Breadcrumb,
    Hamburger
  },
  data() {
    return {
      form: {
        new_pw1: '',
        new_pw2: '',
        old_pw: ''
      },
      dialogFormVisible: false
    }
  },
  computed: {
    ...mapGetters([
      'sidebar',
      'avatar'
    ])
  },
  methods: {
    modify() {
      this.$refs.form.validate(valid => {
        if (valid) {
          resetPass(this.form).then(() => {
            this.$message({
              message: "修改密码成功",
              type: 'success'
            })
            this.dialogFormVisible = false
          })
        }
      })
    },
    resetForm() {
      this.form = {
        new_pw1: '',
        old_pw: '',
        new_pw2: ''
      }
    },
    changePass() {
      this.resetForm()
      this.dialogFormVisible = true
    },
    toggleSideBar() {
      this.$store.dispatch('app/toggleSideBar')
    },
    async logout() {
      await this.$store.dispatch('user/logout')
      this.$router.push(`/login?redirect=${this.$route.fullPath}`)
    }
  }
}
</script>

<style lang="scss" scoped>
.navbar {
  height: 50px;
  overflow: hidden;
  position: relative;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);

  .hamburger-container {
    line-height: 46px;
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background 0.3s;
    -webkit-tap-highlight-color: transparent;

    &:hover {
      background: rgba(0, 0, 0, 0.025);
    }
  }

  .breadcrumb-container {
    float: left;
  }

  .right-menu {
    float: right;
    height: 100%;
    line-height: 50px;

    &:focus {
      outline: none;
    }

    .right-menu-item {
      display: inline-block;
      padding: 0 8px;
      height: 100%;
      font-size: 18px;
      color: #5a5e66;
      vertical-align: text-bottom;

      &.hover-effect {
        cursor: pointer;
        transition: background 0.3s;

        &:hover {
          background: rgba(0, 0, 0, 0.025);
        }
      }
    }

    .avatar-container {
      margin-right: 30px;

      .avatar-wrapper {
        margin-top: 5px;
        position: relative;

        .user-avatar {
          cursor: pointer;
          width: 40px;
          height: 40px;
          border-radius: 10px;
        }

        .el-icon-caret-bottom {
          cursor: pointer;
          position: absolute;
          right: -20px;
          top: 25px;
          font-size: 12px;
        }
      }
    }
  }
}
</style>
