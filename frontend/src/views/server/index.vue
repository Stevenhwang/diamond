<template>
  <el-container>
    <!-- 弹窗 -->
    <el-dialog :title="textMap[dialogStatus]"
               :visible.sync="dialogFormVisible">
      <el-form ref="form"
               :model="form"
               :rules="rules">
        <el-form-item label="名称"
                      :label-width="formLabelWidth"
                      prop="name">
          <el-input v-model="form.name"
                    autocomplete="off" />
        </el-form-item>
        <el-form-item label="IP"
                      :label-width="formLabelWidth"
                      prop="ip">
          <el-input v-model="form.ip"
                    autocomplete="off" />
        </el-form-item>
        <el-form-item label="端口"
                      :label-width="formLabelWidth"
                      prop="port">
          <el-input-number v-model="form.port"
                           :min="1"
                           :max="65535"
                           autocomplete="off" />
        </el-form-item>
        <el-form-item label="认证"
                      :label-width="formLabelWidth"
                      prop="credential_id">
          <el-select v-model="form.credential_id"
                     placeholder="请选择认证">
            <el-option v-for="item in credentials"
                       :key="item.id"
                       :label="item.name"
                       :value="item.id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="实例类型"
                      :label-width="formLabelWidth"
                      prop="instance_type">
          <el-input v-model="form.instance_type"
                    autocomplete="off" />
        </el-form-item>
        <el-form-item label="备注"
                      :label-width="formLabelWidth"
                      prop="remark">
          <el-input v-model="form.remark"
                    type="textarea"
                    autocomplete="off" />
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
    </el-header>
    <el-main>
      <el-table v-loading="listLoading"
                :data="tableData"
                style="width: 100%"
                :row-style="{height:'35px'}"
                :cell-style="{padding:'0 0'}">
        <el-table-column prop="name"
                         label="名称" />
        <el-table-column prop="ip"
                         width="160"
                         label="IP" />
        <el-table-column prop="instance_type"
                         label="实例类型" />
        <el-table-column prop="specifications"
                         show-overflow-tooltip
                         label="实例配置" />
        <el-table-column prop="remark"
                         show-overflow-tooltip
                         label="备注" />
        <!-- <el-table-column prop="created_at"
                         show-overflow-tooltip
                         label="创建时间">
          <template slot-scope="scope">
            {{ parseTime(new Date(scope.row.created_at)) }}
          </template>
        </el-table-column>
        <el-table-column prop="updated_at"
                         show-overflow-tooltip
                         label="更新时间">
          <template slot-scope="scope">
            {{ parseTime(new Date(scope.row.updated_at)) }}
          </template>
        </el-table-column> -->
        <el-table-column label="操作"
                         width="220">
          <template slot-scope="scope">
            <el-button size="mini"
                       type="success"
                       plain
                       @click="handleConn(scope.row)">
              连接
            </el-button>
            <el-button size="mini"
                       type="primary"
                       plain
                       @click="handleEdit(scope.row)">
              编辑
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
import { getCredentials } from '@/api/credential'
import { getServers, createServer, updateServer, deleteServer } from '@/api/server'
import { parseTime } from '@/utils/index'
import Pagination from '@/components/Pagination'

export default {
  components: { Pagination },
  data() {
    return {
      credentials: [],
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
        update: '更新服务器',
        create: '新增服务器'
      },
      form: {
        id: "",
        name: "",
        ip: "",
        port: 22,
        credential_id: "",
        instance_type: "",
        remark: ""
      },
      formLabelWidth: '100px',
      dialogFormVisible: false,
    }
  },
  created() {
    this.getData()
  },
  methods: {
    handleConn(row) {
      let routeUrl = this.$router.resolve({
        path: '/terminal',
        query: {
          id: row.id
        }
      });
      window.open(routeUrl.href, '_blank');
    },
    getCres() {
      getCredentials({ page: 1, limit: 100 }).then(resp => {
        this.credentials = resp.data
      })
    },
    getData() {
      this.listLoading = true
      getServers(this.listQuery).then(resp => {
        this.tableData = resp.data
        this.total = resp.total
        this.listLoading = false
      })
    },
    resetForm() {
      this.form = {
        id: "",
        name: "",
        ip: "",
        port: 22,
        credential_id: "",
        instance_type: "",
        remark: ""
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
          createServer(data).then(resp => {
            console.log(resp)
            this.dialogFormVisible = false
            this.$message.success('新增服务器成功')
            this.getData()
          })
        }
      })
    },
    handleEdit(row) {
      this.getCres()
      this.resetForm()
      this.form.id = row.id
      this.form.name = row.name
      this.form.ip = row.ip
      this.form.port = row.port
      this.form.credential_id = row.credential_id
      this.form.instance_type = row.instance_type
      this.form.remark = row.remark
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
          updateServer(this.form.id, data).then(resp => {
            console.log(resp)
            this.dialogFormVisible = false
            this.$message.success('更新服务器成功');
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
        deleteServer(row.id).then(resp => {
          this.dialogFormVisible = false
          console.log(resp)
          this.$message.success('删除服务器成功');
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
