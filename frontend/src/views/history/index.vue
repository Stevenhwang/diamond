<template>
  <el-container>
    <el-dialog :visible.sync="contentVisible"
               fullscreen>
      <div id="demo"
           style="white-space: pre-line;">{{ content }}</div>
    </el-dialog>
    <el-main>
      <el-table v-loading="listLoading"
                :data="tableData"
                style="width: 100%"
                :row-style="{height:'35px'}"
                :cell-style="{padding:'0 0'}">
        <el-table-column prop="task_name"
                         show-overflow-tooltip
                         label="任务名称" />
        <el-table-column prop="user"
                         label="用户" />
        <el-table-column prop="from_ip"
                         label="来源IP" />
        <el-table-column prop="success"
                         label="状态">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.success"
                    type="success">成功</el-tag>
            <el-tag v-else
                    type="danger">失败</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="content"
                         width="150"
                         label="执行结果">
          <template slot-scope="scope">
            <el-button slot="reference"
                       size="medium"
                       @click="handleClick(scope.row)">
              查看
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="created_at"
                         width="170"
                         label="时间">
          <template slot-scope="scope">
            {{ parseTime(new Date(scope.row.created_at)) }}
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
import { getTaskHist } from '@/api/task'
import { parseTime } from '@/utils/index'
import Pagination from '@/components/Pagination'

export default {
  components: { Pagination },
  data() {
    return {
      contentVisible: false,
      content: "",
      total: 0,
      listQuery: {
        page: 1,
        limit: 15
      },
      parseTime: parseTime,
      tableData: [],
      listLoading: false,
      formLabelWidth: '100px',
      dialogFormVisible: false,
    }
  },
  created() {
    this.getData()
  },
  methods: {
    handleClick(row) {
      this.contentVisible = true
      this.content = row.content
    },
    getData() {
      this.listLoading = true
      getTaskHist(this.listQuery).then(resp => {
        this.tableData = resp.data
        this.total = resp.total
        this.listLoading = false
      })
    }
  },
}
</script>
