<template>
  <el-container>
    <!-- 弹窗 -->
    <el-dialog :visible.sync="recordVisible"
               @closed="handleClose">
      <div id="demo"></div>
    </el-dialog>
    <!-- 弹窗 -->
    <!-- <el-header style="margin-top: 5px"
               height="45px">
      <el-select v-model="searchKey"
                 clearable
                 placeholder="请选择搜索项目"
                 @clear="clearSearchKey"
                 @change="handleChange(searchKey)">
        <el-option v-for="item in options"
                   :key="item.value"
                   :label="item.label"
                   :value="item.value" />
      </el-select>
      <el-input v-if="!isDate"
                v-model="searchValue"
                style="width:350px;"
                clearable
                placeholder="请输入搜索内容"
                @clear="clearSearchValue"
                @keyup.enter.native="handleFilter" />
      <template v-if="isDate">
        <el-date-picker v-model="searchDate"
                        type="datetimerange"
                        range-separator="至"
                        start-placeholder="开始日期"
                        end-placeholder="结束日期"
                        @change="changeDate" />
      </template>
      <el-button class="filter-item"
                 type="primary"
                 icon="el-icon-search"
                 @click="handleFilter">
        搜索
      </el-button>
    </el-header> -->
    <el-main style="padding-top: 5px;padding-bottom: 2px">
      <el-table v-loading="listLoading"
                :data="tableData"
                style="width: 100%"
                :row-style="{height:'35px'}"
                :cell-style="{padding:'0 0'}">
        <el-table-column prop="user"
                         label="用户"
                         width="150" />
        <el-table-column prop="ip"
                         width="160"
                         label="服务器IP" />
        <el-table-column prop="from_ip"
                         width="160"
                         label="来源IP" />
        <el-table-column prop="file"
                         label="记录文件">
          <template slot-scope="scope">
            {{ scope.row.file.replace('./records/', '') }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at"
                         width="170"
                         label="时间">
          <template slot-scope="scope">
            {{ parseTime(new Date(scope.row.created_at)) }}
          </template>
        </el-table-column>
        <el-table-column label="操作"
                         width="200">
          <template slot-scope="scope">
            <el-button type="success"
                       plain
                       icon="el-icon-video-play"
                       size="small"
                       @click="playRecord(scope.row)">
              播放
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
                  @pagination="getList" />
    </el-footer>
  </el-container>
</template>

<script>
import { getRecords } from '@/api/server'
import Pagination from '@/components/Pagination'
import { parseTime } from '@/utils'
import * as AsciinemaPlayer from 'asciinema-player'
import 'asciinema-player/dist/bundle/asciinema-player.css'

export default {
  components: { Pagination },
  data() {
    return {
      player: null,
      recordVisible: false,
      listLoading: false,
      total: 0,
      listQuery: {
        page: 1,
        limit: 15
      },
      parseTime: parseTime,
      formLabelWidth: '100px',
      tableData: []
    }
  },
  created() {
    this.getList()
  },
  methods: {
    handleClose() {
      if (this.player) {
        this.player.dispose()
        this.player = null
      }
    },
    playRecord(row) {
      this.recordVisible = true
      this.$nextTick(() => {
        this.player = AsciinemaPlayer.create(row.file, document.getElementById('demo'), { speed: 2 })
      })
    },
    getList() {
      this.listLoading = true
      getRecords(this.listQuery).then(response => {
        this.tableData = response.data
        this.total = response.total
        this.listLoading = false
      })
    },
  }
}
</script>
