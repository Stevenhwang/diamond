<template>
  <el-container>
    <!-- 弹窗 -->
    <el-dialog :title="textMap[dialogStatus]"
               :visible.sync="dialogFormVisible">
      <el-form ref="form"
               :model="form"
               :rules="rules">
        <el-form-item label="用户名"
                      :label-width="formLabelWidth"
                      prop="username">
          <el-input v-model="form.username"
                    autocomplete="off" />
        </el-form-item>
        <el-form-item label="密码"
                      :label-width="formLabelWidth"
                      prop="password">
          <el-input v-model="form.password"
                    autocomplete="off" />
        </el-form-item>
        <el-form-item label="公钥"
                      :label-width="formLabelWidth"
                      prop="publickey">
          <el-input v-model="form.publickey"
                    type="textarea"
                    :rows="4"
                    autocomplete="off" />
        </el-form-item>
        <el-form-item label="菜单"
                      :label-width="formLabelWidth"
                      prop="menus">
          <el-select v-model="form.menus"
                     multiple
                     placeholder="请选择">
            <el-option v-for="item in menuOptions"
                       :key="item.value"
                       :label="item.label"
                       :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="激活"
                      :label-width="formLabelWidth">
          <el-switch v-model="form.is_active" />
        </el-form-item>
      </el-form>
      <div slot="footer"
           class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          取 消
        </el-button>
        <el-button type="primary"
                   @click="dialogStatus==='create'?createData():updateData()">
          确 定
        </el-button>
      </div>
    </el-dialog>
    <!-- 弹窗 -->
    <!-- 权限分配 -->
    <el-dialog title="权限分配"
               fullscreen
               :visible.sync="dialogTableVisible">
      <el-table ref="multipleTable"
                :data="gridData"
                @selection-change="handleSelectionChange">
        <el-table-column v-if="selectUser !== 'admin'"
                         type="selection"
                         width="50">
        </el-table-column>
        <el-table-column property="name"
                         label="姓名">
        </el-table-column>
        <el-table-column property="method"
                         label="请求方法">
        </el-table-column>
        <el-table-column property="url"
                         label="请求URL">
        </el-table-column>
      </el-table>
      <div v-if="selectUser !== 'admin'"
           slot="footer"
           class="dialog-footer">
        <el-button @click="dialogTableVisible = false">
          取 消
        </el-button>
        <el-button type="primary"
                   @click="assignPerm()">
          确 定
        </el-button>
      </div>
    </el-dialog>
    <!-- 权限分配 -->

    <!-- 服务器分配 -->
    <el-dialog title="服务器分配"
               fullscreen
               :visible.sync="serverTableVisible">
      <el-table ref="serverTable"
                :data="serverData"
                @selection-change="handleServerChange">
        <el-table-column v-if="selectUser !== 'admin'"
                         type="selection"
                         width="50">
        </el-table-column>
        <el-table-column property="name"
                         label="名称">
        </el-table-column>
        <el-table-column property="ip"
                         label="IP">
        </el-table-column>
        <el-table-column property="remark"
                         show-overflow-tooltip
                         label="备注">
        </el-table-column>
      </el-table>
      <div v-if="selectUser !== 'admin'"
           slot="footer"
           class="dialog-footer">
        <el-button @click="serverTableVisible = false">
          取 消
        </el-button>
        <el-button type="primary"
                   @click="assignServer()">
          确 定
        </el-button>
      </div>
    </el-dialog>
    <!-- 服务器分配 -->

    <el-header style="margin-top: 5px"
               height="30px">
      <!-- <el-select v-model="searchKey"
                 size="small"
                 clearable
                 placeholder="请选择搜索项目"
                 @clear="clearSearchKey">
        <el-option v-for="item in options"
                   :key="item.value"
                   :label="item.label"
                   :value="item.value" />
      </el-select>
      <el-input v-model="searchValue"
                size="small"
                style="width:350px;"
                clearable
                placeholder="请输入搜索内容"
                @clear="clearSearchValue"
                @keyup.enter.native="handleFilter" />
      <el-button class="filter-item"
                 type="primary"
                 icon="el-icon-search"
                 size="small"
                 @click="handleFilter">
        搜索
      </el-button> -->
      <el-button type="primary"
                 size="small"
                 icon="el-icon-edit"
                 @click="handleCreate()">新增
      </el-button>
      <el-button type="success"
                 size="small"
                 icon="el-icon-refresh"
                 @click="handleSync()">同步权限
      </el-button>
    </el-header>
    <el-main>
      <el-table v-loading="listLoading"
                :data="tableData"
                style="width: 100%"
                :row-style="{height:'35px'}"
                :cell-style="{padding:'0 0'}">
        <el-table-column prop="username"
                         label="用户">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.is_active"
                    type="success">
              {{ scope.row.username }}
            </el-tag>
            <el-tag v-else
                    type="danger">
              {{ scope.row.username }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_login_ip"
                         label="登录IP" />
        <el-table-column prop="last_login_time"
                         label="登录时间">
          <template slot-scope="scope">
            {{ parseTime(new Date(scope.row.last_login_time)) }}
          </template>
        </el-table-column>
        <!-- <el-table-column prop="created_at"
                         width="160"
                         label="创建时间">
          <template slot-scope="scope">
            {{ parseTime(new Date(scope.row.created_at)) }}
          </template>
        </el-table-column>
        <el-table-column prop="updated_at"
                         width="160"
                         label="更新时间">
          <template slot-scope="scope">
            {{ parseTime(new Date(scope.row.updated_at)) }}
          </template>
        </el-table-column> -->
        <el-table-column label="操作"
                         width="300">
          <template slot-scope="scope">
            <el-button size="mini"
                       type="primary"
                       plain
                       @click="handleEdit(scope.row)">
              编辑
            </el-button>
            <el-button size="mini"
                       type="success"
                       plain
                       @click="handlePerm(scope.row)">
              权限
            </el-button>
            <el-button size="mini"
                       type="info"
                       plain
                       @click="handleServer(scope.row)">
              服务器
            </el-button>
            <el-button size="mini"
                       type="danger"
                       plain
                       @click="handleDelete(scope.row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-main>
    <el-footer>
      <pagination v-show="total>0"
                  :total="total"
                  :page.sync="listQuery.page"
                  :limit.sync="listQuery.limit"
                  @pagination="getData" />
    </el-footer>
  </el-container>
</template>

<script>
import {
  getUsers, createUser, updateUser, deleteUser, getUserPerms,
  assignUserPerm, getUserServers, assignUserServer, syncPerms
} from '@/api/user'
import { parseTime } from '@/utils/index'
import Pagination from '@/components/Pagination'

export default {
  components: { Pagination },
  data() {
    return {
      serverTableVisible: false,
      serverData: [],
      serverSelect: [],
      dialogTableVisible: false,
      gridData: [],
      multipleSelection: [],
      selectUser: "",
      menuOptions: [{
        value: 'server',
        label: '服务器'
      }, {
        value: 'credential',
        label: '认证'
      }, {
        value: 'record',
        label: 'ssh记录'
      }, {
        value: 'user',
        label: '用户'
      }],
      total: 0,
      listQuery: {
        page: 1,
        limit: 15
      },
      parseTime: parseTime,
      tableData: [],
      listLoading: false,
      rules: {},
      dialogStatus: '',
      textMap: {
        update: '更新用户',
        create: '新增用户'
      },
      form: {
        id: "",
        username: "",
        password: "",
        publickey: "",
        menus: [],
        is_active: true
      },
      formLabelWidth: '100px',
      dialogFormVisible: false,
    }
  },
  created() {
    this.getData()
  },
  methods: {
    handleSync() {
      syncPerms().then(() => {
        this.$message.success('同步权限成功')
      })
    },
    handleServerChange(val) {
      this.serverSelect = val
    },
    handleSelectionChange(val) {
      this.multipleSelection = val
    },
    handleServer(row) {
      this.serverTableVisible = true
      this.selectUser = row.username
      getUserServers({ id: row.id }).then(resp => {
        this.serverData = resp.data
        this.$nextTick(() => {
          this.serverData.forEach((e, i) => {
            this.$refs.serverTable.toggleRowSelection(this.serverData[i], e.check)
          })
        })
      })
    },
    assignServer() {
      let servers = []
      this.serverSelect.forEach(e => {
        servers.push(e.id)
      })
      assignUserServer({ username: this.selectUser, servers: servers }).then(() => {
        this.$message.success('分配服务器成功')
        this.serverTableVisible = false
      })
    },
    handlePerm(row) {
      this.dialogTableVisible = true
      this.selectUser = row.username
      getUserPerms({ id: row.id }).then(resp => {
        this.gridData = resp.data
        this.$nextTick(() => {
          this.gridData.forEach((e, i) => {
            this.$refs.multipleTable.toggleRowSelection(this.gridData[i], e.check)
          })
        })
      })
    },
    assignPerm() {
      let perms = []
      this.multipleSelection.forEach(e => {
        perms.push(e.name)
      })
      assignUserPerm({ username: this.selectUser, perms: perms }).then(() => {
        this.$message.success('分配权限成功')
        this.dialogTableVisible = false
      })
    },
    getData() {
      this.listLoading = true
      getUsers(this.listQuery).then(resp => {
        this.tableData = resp.data
        this.total = resp.total
        this.listLoading = false
      })
    },
    resetForm() {
      this.form = {
        id: "",
        username: "",
        password: "",
        publickey: "",
        menus: [],
        is_active: true
      }
    },
    handleCreate() {
      this.resetForm()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['form'].clearValidate()
      })
    },
    createData() {
      this.$refs['form'].validate((valid) => {
        if (valid) {
          const data = Object.assign({}, this.form)
          delete data.id
          createUser(data).then(() => {
            this.dialogFormVisible = false
            this.$message.success('新增用户成功')
            this.getData()
          })
        }
      })
    },
    handleEdit(row) {
      this.resetForm()
      this.form.id = row.id
      this.form.username = row.username
      this.form.password = row.password
      this.form.publickey = row.publickey
      this.form.menus = row.menus
      this.form.is_active = row.is_active
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['form'].clearValidate()
      })
    },
    updateData() {
      this.$refs['form'].validate((valid) => {
        if (valid) {
          const data = Object.assign({}, this.form)
          delete data.id
          updateUser(this.form.id, data).then(() => {
            this.dialogFormVisible = false
            this.$message.success('更新用户成功');
            this.getData()
          })
        }
      })
    },
    handleDelete(row) {
      this.$confirm('确认删除?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteUser(row.id).then(() => {
          this.dialogFormVisible = false
          this.$message.success('删除用户成功');
          this.getData()
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        });
      });
    }
  },
}
</script>
