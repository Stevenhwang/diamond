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
        <el-form-item label="目标"
                      :label-width="formLabelWidth"
                      prop="target">
          <el-select v-model="form.target"
                     style="width:100%"
                     filterable
                     allow-create
                     placeholder="请选择服务器或填写服务器分组">
            <el-option v-for="item in servers"
                       :key="item.id"
                       :label="item.name + '(' + item.ip + ')'"
                       :value="item.ip">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="脚本"
                      :label-width="formLabelWidth"
                      prop="script_id">
          <el-select v-model="form.script_id"
                     style="width:100%"
                     filterable
                     placeholder="请选择要执行的脚本">
            <el-option v-for="item in scripts"
                       :key="item.id"
                       :label="item.name"
                       :value="item.id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="参数"
                      :label-width="formLabelWidth"
                      prop="args">
          <el-input v-model="form.args"
                    placeholder="请填写脚本执行参数，多个参数空格分开"
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
      <el-input v-model="listQuery.query"
                size="small"
                style="width:350px;"
                prefix-icon="el-icon-search"
                clearable
                placeholder="请输入搜索内容，支持名称和命令"
                @input="changeSearch" />&nbsp;
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
                         width="300"
                         show-overflow-tooltip
                         label="名称" />
        <el-table-column prop="target"
                         show-overflow-tooltip
                         label="目标" />
        <el-table-column label="操作"
                         width="220">
          <template slot-scope="scope">
            <el-button size="mini"
                       type="success"
                       plain
                       @click="handleInvoke(scope.row)">
              执行
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
import { getTasks, createTask, updateTask, deleteTask, invokeTask } from '@/api/task'
import { getScripts } from '@/api/script'
import { getServers } from '@/api/server'
import { parseTime } from '@/utils/index'
import Pagination from '@/components/Pagination'

export default {
  components: { Pagination },
  data() {
    return {
      scripts: [],
      servers: [],
      total: 0,
      listQuery: {
        page: 1,
        limit: 15,
        query: ""
      },
      parseTime: parseTime,
      tableData: [],
      listLoading: false,
      rules: {},
      dialogStatus: '',
      textMap: {
        update: '更新任务',
        create: '新增任务'
      },
      form: {
        id: "",
        name: "",
        target: "",
        script_id: "",
        args: ""
      },
      formLabelWidth: '100px',
      dialogFormVisible: false,
    }
  },
  created() {
    this.getData()
  },
  methods: {
    getScrs() {
      getScripts({ page: 1, limit: 100 }).then(resp => {
        this.scripts = resp.data
      })
    },
    getSers() {
      getServers({ page: 1, limit: 100 }).then(resp => {
        this.servers = resp.data
      })
    },
    changeSearch() {
      this.listQuery.page = 1
      this.getData()
    },
    handleInvoke(row) {
      this.$confirm('确认执行?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        invokeTask(row.id).then(() => {
          this.dialogFormVisible = false
          this.$message.success('执行任务成功')
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消执行'
        });
      });
    },
    getData() {
      this.listLoading = true
      getTasks(this.listQuery).then(resp => {
        this.tableData = resp.data
        this.total = resp.total
        this.listLoading = false
      })
    },
    resetForm() {
      this.form = {
        id: "",
        name: "",
        target: "",
        script_id: "",
        args: ""
      }
    },
    handleCreate() {
      this.getScrs()
      this.getSers()
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
          createTask(data).then(() => {
            this.dialogFormVisible = false
            this.$message.success('新增任务成功')
            this.getData()
          })
        }
      })
    },
    handleEdit(row) {
      this.getScrs()
      this.getSers()
      this.resetForm()
      this.form.id = row.id
      this.form.name = row.name
      this.form.target = row.target
      this.form.script_id = row.script_id
      this.form.args = row.args
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
          updateTask(this.form.id, data).then(() => {
            this.dialogFormVisible = false
            this.$message.success('更新任务成功');
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
        deleteTask(row.id).then(() => {
          this.dialogFormVisible = false
          this.$message.success('删除任务成功');
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
