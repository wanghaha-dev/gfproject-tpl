<template>
    <div>
        <el-button type="primary" round @click="addDialogVisible=!addDialogVisible">添加[[.TypeName]]</el-button>
        <el-table
            :data="dataList"
            style="width: 100%">

            [[range .AllFields]]
           <el-table-column
           label="[[if .Comment]][[.Comment]][[else]][[.Name]][[end]]"
           width="[[if TypeCheck .Type]]200[[else]]50[[end]]">
           <template slot-scope="scope">
               <span style="margin-left: 10px">{{ scope.row.[[.Name]] }}</span>
           </template>
           </el-table-column>
            [[end]]

            <el-table-column label="操作" width="200">
                <template slot-scope="scope">
                    <el-button type="primary" icon="el-icon-edit" circle
                                    @click="handleEdit(scope.$index, scope.row)"></el-button>

                    <span style="margin-left: 5px;">
                        <el-popconfirm
                        confirm-button-text='好的'
                        cancel-button-text='不用了'
                        icon="el-icon-info"
                        icon-color="red"
                        title="这是一段内容确定删除吗？"
                        @confirm="handleDelete(scope.$index, scope.row)"
                        >

                        <el-button type="danger" icon="el-icon-delete" circle slot="reference"></el-button>
                    </el-popconfirm>
                    </span>
                </template>
            </el-table-column>
        </el-table>

        <div class="pagination-view">
            <el-pagination
            background
            layout="sizes, prev, pager, next, total"
            @current-change="changePage"
            :page-size="pageSize"
            @size-change="handleSizeChange"
            :page-sizes="[10, 20, 50, 100]"
            :total="total">
            </el-pagination>
        </div>


        <!-- 添加[[.TypeName]]弹框 -->
        <el-dialog @open="onAddOpen" @close="onAddClose" title="添加[[.TypeName]]" :visible.sync="addDialogVisible">
            <el-form ref="addForm" :model="addFormData" :rules="rules" size="medium" label-width="100px">
                [[range .Fields]]
                <el-form-item label-width="300" label="[[if .Comment]][[.Comment]][[else]][[.Name]][[end]]" prop="[[.Name]]">
                <el-input v-model="addFormData.[[.Name]]" placeholder="请输入[[if .Comment]][[.Comment]][[else]][[.Name]][[end]]" clearable :style="{width: '100%'}">
                </el-input>
                </el-form-item>
                [[end]]
            </el-form>
            <div slot="footer">
                <el-button @click="closeAdd">取消</el-button>
                <el-button type="primary" @click="addSubmit">确定</el-button>
            </div>
        </el-dialog>

        <!-- 编辑[[.TypeName]]弹框 -->
        <el-dialog @open="onEditOpen" @close="onEditClose" title="编辑[[.TypeName]]" :visible.sync="editDialogVisible">
            <el-form ref="editForm" :model="editFormData" :rules="rules" size="medium" label-width="100px">
                [[range .Fields]]
                <el-form-item label-width="300" label="[[if .Comment]][[.Comment]][[else]][[.Name]][[end]]" prop="[[.Name]]">
                <el-input v-model="editFormData.[[.Name]]" placeholder="请输入[[if .Comment]][[.Comment]][[else]][[.Name]][[end]]" clearable :style="{width: '100%'}">
                </el-input>
                </el-form-item>
                [[end]]
            </el-form>
            <div slot="footer">
                <el-button @click="closeEdit">取消</el-button>
                <el-button type="primary" @click="editSubmit">确定</el-button>
            </div>
        </el-dialog>

    </div>
</template>

<script>
import axios from "axios";
export default {
    data() {
        return {
            total: 0,
            currentPage: 1,
            pageSize: null,
            dataList: [],
            addDialogVisible: false,
            editDialogVisible: false,
            addFormData: {
                [[range .Fields]] [[.Name]]: undefined,
                [[end]]
            },
            editFormData: {
                [[range .Fields]] [[.Name]]: undefined,
                [[end]]
            },
            rules: {
            [[range .Fields]] [[if .Null]] [[.Name]]: [], [[else]] [[.Name]]: [{
                required: true,
                message: '请输入[[if .Comment]][[.Comment]][[else]][[.Name]][[end]]',
                trigger: 'blur'
                }],[[end]]
                [[end]]
            },
        }
    },
    mounted() {
        this.getDataList()
    },
    methods: {
        getDataList() {
            axios.get("http://127.0.0.1:8199/[[.RouterName]]/list?page="+this.currentPage+"&pageSize="+this.pageSize+"").then(ret=> {
                if(ret.data.code==0) {
                    this.total = ret.data.total;
                    this.pageSize = ret.data.pageSize;
                    this.dataList = ret.data.data;
                } else {
                    this.$notify({
                        title: ret.data.msg,
                        message: ret.data.msg,
                        type: 'error'
                    });
                }
            })
        },
        handleEdit(index, row) {
            console.log(index, row);
            this.editFormData = row;
            this.editDialogVisible = true;
        },
        handleDelete(index, row) {
            axios.delete("http://127.0.0.1:8199/[[.RouterName]]/"+row.id+"/delete").then(ret=>{
                if(ret.data.code==0) {
                    this.$notify({
                        title: '成功',
                        message: ret.data.msg,
                        type: 'success'
                    });
                    this.getDataList();
                } else {
                    this.$notify({
                        title: ret.data.msg,
                        message: ret.data.msg,
                        type: 'error'
                    });
                }
            })
            console.log(index, row);
        },

        //
        onAddOpen() {},
        onEditOpen() {},
        onAddClose() {
            this.$refs['addForm'].resetFields()
        },
        onEditClose() {
            this.$refs['editForm'].resetFields()
        },
        closeAdd() {
            this.addDialogVisible = false;
        },
        closeEdit() {
            this.editDialogVisible = false;
        },
        addSubmit() {
            this.$refs['addForm'].validate((valid) => {
                if (valid) {
                    // alert('submit!');
                    axios.post("http://127.0.0.1:8199/[[.RouterName]]/add", this.addFormData).then(ret=>{
                        if(ret.data.code==0) {
                            this.$notify({
                                title: '成功',
                                message: ret.data.msg,
                                type: 'success'
                            });
                            this.closeAdd()
                            this.addDialogVisible = false;
                            this.getDataList();
                        } else {
                            this.$notify({
                                title: '失败',
                                message: ret.data.msg,
                                type: 'error'
                            });
                        }
                    })
                } else {
                    console.log('error submit!!');
                    return false;
                }
            });
        },
        editSubmit() {
            this.$refs['editForm'].validate((valid) => {
                if (valid) {
                    axios.post("http://127.0.0.1:8199/[[.RouterName]]/"+this.editFormData.id+"/update", this.editFormData).then(ret=>{
                        if(ret.data.code==0) {
                            this.$notify({
                                title: '成功',
                                message: ret.data.msg,
                                type: 'success'
                            });
                            this.onEditClose()
                            this.editDialogVisible = false;
                            this.getDataList();
                        } else {
                            this.$notify({
                                title: '失败',
                                message: ret.data.msg,
                                type: 'error'
                            });
                        }
                    })
                } else {
                    console.log('error submit!!');
                    return false;
                }
            });
        },
        changePage(currentPage) {
            this.currentPage = currentPage;
            this.getDataList();
        },
        handleSizeChange(size) {
            this.pageSize = size;
            this.getDataList();
        }
    }
}
</script>

<style scoped>
.pagination-view {
    margin: 20px;
}
</style>