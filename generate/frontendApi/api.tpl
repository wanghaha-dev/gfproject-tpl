import request from '@/utils/request';

/**
 * 分页查询{{.PackageNameNotes}}
 * @param params 查询参数
 */
export async function page{{.PackageFuncName}}(params) {
    const res = await request.get('/{{.PackageName}}/list', {
        params
    });
    if (res.data.code === 0) {
        return res.data.data;
    }
    return Promise.reject(new Error(res.data.message));
}

/**
 * 批量修改{{.PackageNameNotes}}
 * @param data 参数
 */
export async function update{{.PackageFuncName}}(data) {
    const res = await request.put('/{{.PackageName}}/updateBatch', data);
    if (res.data.code === 0) {
        return res.data.message;
    }
    return Promise.reject(new Error(res.data.message));
}

/**
 * 修改单个{{.PackageNameNotes}}
 * @param data 参数
 */
export async function updateOne{{.PackageFuncName}}(data) {
    const res = await request.put('/{{.PackageName}}/updateOne', data);
    if (res.data.code === 0) {
        return res.data.message;
    }
    return Promise.reject(new Error(res.data.message));
}

/**
 * 批量删除{{.PackageNameNotes}}
 * @param data 参数
 */
export async function deleteBatch{{.PackageFuncName}}(data) {
    const res = await request.delete('/{{.PackageName}}/deleteBatch', data);
    if (res.data.code === 0) {
        return res.data.message;
    }
    return Promise.reject(new Error(res.data.message));
}

/**
 * 删除单个{{.PackageNameNotes}}
 * @param data 参数
 */
export async function deleteOne{{.PackageFuncName}}(data) {
    const res = await request.delete('/{{.PackageName}}/deleteOne', data);
    if (res.data.code === 0) {
        return res.data.message;
    }
    return Promise.reject(new Error(res.data.message));
}
