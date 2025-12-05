'use client'

import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import api from '@/lib/api'

interface Barang {
  id: number
  kode_barang: string
  nama_barang: string
  kategori: string
  satuan: string
  harga_beli: number
  harga_jual: number
}

export default function BarangPage() {
  const [barangs, setBarangs] = useState<Barang[]>([])
  const [loading, setLoading] = useState(false)
  const [showModal, setShowModal] = useState(false)
  const [editId, setEditId] = useState<number | null>(null)
  const [search, setSearch] = useState('')
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const limit = 10

  const [formData, setFormData] = useState({
    kode_barang: '',
    nama_barang: '',
    kategori: '',
    satuan: '',
    harga_beli: 0,
    harga_jual: 0
  })

  useEffect(() => {
    fetchBarangs()
  }, [page, search])

  const fetchBarangs = async () => {
    setLoading(true)
    try {
      const response = await api.get(`/barang?search=${search}&page=${page}&limit=${limit}`)
      setBarangs(response.data.data || [])
      setTotal(response.data.meta?.total || 0)
    } catch (error) {
      console.error('Error fetching barang:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      if (editId) {
        await api.put(`/barang/${editId}`, formData)
        toast.success('Barang updated successfully!')
      } else {
        await api.post('/barang', formData)
        toast.success('Barang created successfully!')
      }
      setShowModal(false)
      resetForm()
      fetchBarangs()
    } catch (error: any) {
      toast.error(error.response?.data?.message || 'Operation failed')
    } finally {
      setLoading(false)
    }
  }

  const handleEdit = (barang: Barang) => {
    setEditId(barang.id)
    setFormData({
      kode_barang: barang.kode_barang,
      nama_barang: barang.nama_barang,
      kategori: barang.kategori,
      satuan: barang.satuan,
      harga_beli: barang.harga_beli,
      harga_jual: barang.harga_jual
    })
    setShowModal(true)
  }

  const handleDelete = async (id: number) => {
    toast((t) => (
      <div className="flex flex-col gap-3">
        <div className="font-semibold">Konfirmasi Hapus</div>
        <div>Apakah Anda yakin ingin menghapus barang ini?</div>
        <div className="flex gap-2 justify-end">
          <button
            onClick={() => toast.dismiss(t.id)}
            className="px-4 py-2 bg-gray-200 hover:bg-gray-300 rounded text-sm font-medium"
          >
            Batal
          </button>
          <button
            onClick={async () => {
              toast.dismiss(t.id)
              try {
                await api.delete(`/barang/${id}`)
                toast.success('Barang berhasil dihapus!')
                fetchBarangs()
              } catch (error: any) {
                toast.error(error.response?.data?.message || 'Gagal menghapus barang')
              }
            }}
            className="px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded text-sm font-medium"
          >
            Hapus
          </button>
        </div>
      </div>
    ), {
      duration: Infinity,
      style: {
        minWidth: '320px'
      }
    })
  }

  const resetForm = () => {
    setEditId(null)
    setFormData({
      kode_barang: '',
      nama_barang: '',
      kategori: '',
      satuan: '',
      harga_beli: 0,
      harga_jual: 0
    })
  }

  const totalPages = Math.ceil(total / limit)

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Master Barang</h1>
        <button
          onClick={() => {
            resetForm()
            setShowModal(true)
          }}
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          + Tambah Barang
        </button>
      </div>

      {/* Search */}
      <div className="mb-4">
        <input
          type="text"
          placeholder="Search by kode or nama barang..."
          value={search}
          onChange={(e) => {
            setSearch(e.target.value)
            setPage(1)
          }}
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
        />
      </div>

      {/* Table */}
      <div className="bg-white shadow-md rounded-lg overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Kode</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nama Barang</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Kategori</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Satuan</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Harga Beli</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Harga Jual</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {loading ? (
              <tr>
                <td colSpan={7} className="px-6 py-4 text-center">Loading...</td>
              </tr>
            ) : barangs.length === 0 ? (
              <tr>
                <td colSpan={7} className="px-6 py-4 text-center">No data found</td>
              </tr>
            ) : (
              barangs.map((barang) => (
                <tr key={barang.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{barang.kode_barang}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{barang.nama_barang}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{barang.kategori}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{barang.satuan}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">Rp {barang.harga_beli.toLocaleString()}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">Rp {barang.harga_jual.toLocaleString()}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <button
                      onClick={() => handleEdit(barang)}
                      className="text-indigo-600 hover:text-indigo-900 mr-3"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDelete(barang.id)}
                      className="text-red-600 hover:text-red-900"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      <div className="mt-4 flex justify-between items-center">
        <div className="text-sm text-gray-700">
          Showing {barangs.length} of {total} results
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => setPage(page - 1)}
            disabled={page === 1}
            className="px-4 py-2 border rounded disabled:opacity-50"
          >
            Previous
          </button>
          <span className="px-4 py-2">
            Page {page} of {totalPages}
          </span>
          <button
            onClick={() => setPage(page + 1)}
            disabled={page === totalPages}
            className="px-4 py-2 border rounded disabled:opacity-50"
          >
            Next
          </button>
        </div>
      </div>

      {/* Modal */}
      {showModal && (
        <div className="fixed z-10 inset-0 overflow-y-auto">
          <div className="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={() => setShowModal(false)}></div>

            <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
              <form onSubmit={handleSubmit}>
                <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                  <h3 className="text-lg font-medium text-gray-900 mb-4">
                    {editId ? 'Edit Barang' : 'Tambah Barang'}
                  </h3>

                  <div className="mb-4">
                    <label className="block text-gray-700 text-sm font-bold mb-2">Kode Barang</label>
                    <input
                      type="text"
                      value={formData.kode_barang}
                      onChange={(e) => setFormData({ ...formData, kode_barang: e.target.value })}
                      className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700"
                      required
                      disabled={!!editId}
                    />
                  </div>

                  <div className="mb-4">
                    <label className="block text-gray-700 text-sm font-bold mb-2">Nama Barang</label>
                    <input
                      type="text"
                      value={formData.nama_barang}
                      onChange={(e) => setFormData({ ...formData, nama_barang: e.target.value })}
                      className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700"
                      required
                    />
                  </div>

                  <div className="mb-4">
                    <label className="block text-gray-700 text-sm font-bold mb-2">Kategori</label>
                    <input
                      type="text"
                      value={formData.kategori}
                      onChange={(e) => setFormData({ ...formData, kategori: e.target.value })}
                      className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700"
                    />
                  </div>

                  <div className="mb-4">
                    <label className="block text-gray-700 text-sm font-bold mb-2">Satuan</label>
                    <input
                      type="text"
                      value={formData.satuan}
                      onChange={(e) => setFormData({ ...formData, satuan: e.target.value })}
                      className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700"
                      required
                    />
                  </div>

                  <div className="mb-4">
                    <label className="block text-gray-700 text-sm font-bold mb-2">Harga Beli</label>
                    <input
                      type="number"
                      value={formData.harga_beli}
                      onChange={(e) => setFormData({ ...formData, harga_beli: Number(e.target.value) })}
                      className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700"
                    />
                  </div>

                  <div className="mb-4">
                    <label className="block text-gray-700 text-sm font-bold mb-2">Harga Jual</label>
                    <input
                      type="number"
                      value={formData.harga_jual}
                      onChange={(e) => setFormData({ ...formData, harga_jual: Number(e.target.value) })}
                      className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700"
                    />
                  </div>
                </div>

                <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
                  <button
                    type="submit"
                    disabled={loading}
                    className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
                  >
                    {loading ? 'Saving...' : 'Save'}
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      setShowModal(false)
                      resetForm()
                    }}
                    className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
                  >
                    Cancel
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
