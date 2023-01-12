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
        <el-form-item label="认证类型"
                      :label-width="formLabelWidth"
                      prop="auth_type">
          <el-select v-model="form.auth_type"
                     placeholder="选择认证类型">
            <el-option v-for="item in options"
                       :key="item.value"
                       :label="item.label"
                       :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="认证用户"
                      :label-width="formLabelWidth"
                      prop="auth_user">
          <el-input v-model="form.auth_user"
                    autocomplete="off" />
        </el-form-item>
        <el-form-item label="认证内容"
                      :label-width="formLabelWidth"
                      prop="auth_content">
          <el-input v-model="form.auth_content"
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
        <el-table-column prop="auth_type"
                         label="认证类型">
          <template slot-scope="scope">
            {{ scope.row.auth_type === 1? "密码认证": "私钥认证" }}
          </template>
        </el-table-column>
        <el-table-column prop="auth_user"
                         label="认证用户" />
        <el-table-column prop="created_at"
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
        </el-table-column>
        <el-table-column label="操作"
                         width="220">
          <template slot-scope="scope">
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
import { getCredentials, createCredential, updateCredential, deleteCredential } from '@/api/credential'
import { parseTime } from '@/utils/index'
import Pagination from '@/components/Pagination'

export default {
  components: { Pagination },
  data() {
    return {
      options: [{
        value: 1,
        label: '密码认证'
      }, {
        value: 2,
        label: '私钥认证'
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
        update: '更新认证',
        create: '新增认证'
      },
      form: {
        id: "",
        name: "",
        auth_type: 1,
        auth_user: "root",
        auth_content: ""
      },
      formLabelWidth: '100px',
      dialogFormVisible: false,
    }
  },
  created() {
    this.getData()
  },
  methods: {
    getData() {
      this.listLoading = true
      getCredentials(this.listQuery).then(resp => {
        this.tableData = resp.data
        this.total = resp.total
        this.listLoading = false
      })
    },
    resetForm() {
      this.form = {
        id: "",
        name: "",
        auth_type: 1,
        auth_user: "root",
        auth_content: ""
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
          createCredential(data).then(() => {
            this.dialogFormVisible = false
            this.$message.success('新增认证成功')
            this.getData()
          })
        }
      })
    },
    handleEdit(row) {
      this.resetForm()
      this.form.id = row.id
      this.form.name = row.name
      this.form.auth_type = row.auth_type
      this.form.auth_user = row.auth_user
      this.form.auth_content = row.auth_content
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
          updateCredential(this.form.id, data).then(() => {
            this.dialogFormVisible = false
            this.$message.success('更新认证成功');
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
        deleteCredential(row.id).then(() => {
          this.dialogFormVisible = false
          this.$message.success('删除认证成功');
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
