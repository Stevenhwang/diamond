<template>
  <el-container>
    <el-header style="margin-top: 5px"
               height="30px">
      <el-input v-model="listQuery.ip"
                size="small"
                style="width:350px;"
                prefix-icon="el-icon-search"
                clearable
                placeholder="请输入搜索IP黑名单"
                @input="changeSearch" />
    </el-header>
    <el-main>
      <el-row>
        <el-col v-for="tag in tableData"
                :key="tag"
                :span="3">
          <el-tag closable
                  @close="handleClose(tag)">
            {{tag}}
          </el-tag>
        </el-col>
      </el-row>
    </el-main>
  </el-container>
</template>

<script>

import { getBanIPs, delBanIP } from '@/api/user'

export default {
  data() {
    return {
      total: 0,
      listQuery: {
        ip: ""
      },
      tableData: [],
      listLoading: false,
      dialogStatus: '',
    }
  },
  created() {
    this.getData()
  },
  methods: {
    changeSearch() {
      this.getData()
    },
    getData() {
      this.listLoading = true
      getBanIPs(this.listQuery).then(resp => {
        this.tableData = resp.data
        this.listLoading = false
      })
    },
    handleClose(tag) {
      this.$confirm('确认删除?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        delBanIP({ ip: tag }).then(() => {
          this.dialogFormVisible = false
          this.$message.success('删除IP成功');
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
