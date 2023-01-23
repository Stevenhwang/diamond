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
        <el-form-item label="时间"
                      :label-width="formLabelWidth"
                      prop="spec">
          <el-input v-model="form.spec"
                    placeholder="支持秒级的crontab表达式(* * * * * *)"
                    autocomplete="off" />
        </el-form-item>
        <el-form-item label="命令"
                      :label-width="formLabelWidth"
                      prop="command">
          <el-input v-model="form.command"
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
        <el-table-column prop="spec"
                         width="100"
                         show-overflow-tooltip
                         label="时间" />
        <el-table-column prop="command"
                         show-overflow-tooltip
                         label="命令" />
        <el-table-column label="操作"
                         width="160">
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
import { getCrons, createCron, updateCron, deleteCron } from '@/api/cron'
import { parseTime } from '@/utils/index'
import Pagination from '@/components/Pagination'

export default {
  components: { Pagination },
  data() {
    return {
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
        update: '更新定时任务',
        create: '新增定时任务'
      },
      form: {
        id: "",
        name: "",
        command: "",
        spec: ""
      },
      formLabelWidth: '100px',
      dialogFormVisible: false,
    }
  },
  created() {
    this.getData()
  },
  methods: {
    changeSearch() {
      this.listQuery.page = 1
      this.getData()
    },
    getData() {
      this.listLoading = true
      getCrons(this.listQuery).then(resp => {
        this.tableData = resp.data
        this.total = resp.total
        this.listLoading = false
      })
    },
    resetForm() {
      this.form = {
        id: "",
        name: "",
        command: "",
        spec: ""
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
          createCron(data).then(() => {
            this.dialogFormVisible = false
            this.$message.success('新增定时任务成功')
            this.getData()
          })
        }
      })
    },
    handleEdit(row) {
      this.resetForm()
      this.form.id = row.id
      this.form.name = row.name
      this.form.command = row.command
      this.form.spec = row.spec
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
          updateCron(this.form.id, data).then(() => {
            this.dialogFormVisible = false
            this.$message.success('更新定时任务成功');
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
        deleteCron(row.id).then(() => {
          this.dialogFormVisible = false
          this.$message.success('删除定时任务成功');
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