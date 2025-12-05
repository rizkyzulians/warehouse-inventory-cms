'use client'

import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import api from '@/lib/api'

interface Pembelian {
  id: number
  no_faktur: string
  tanggal: string
  supplier: string
  total: number
  keterangan: string
}

interface Barang {
  id: number
  kode_barang: string
  nama_barang: string
}

interface DetailItem {
  barang_id: number
  qty: number
  harga: number
}

export default function PembelianPage() {
  const [pembelians, setPembelians] = useState<Pembelian[]>([])
  const [barangs, setBarangs] = useState<Barang[]>([])
  const [loading, setLoading] = useState(false)
  const [showModal, setShowModal] = useState(false)
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const limit = 10

  const [formData, setFormData] = useState({
    no_faktur: '',
    tanggal: '',
    supplier: '',
    keterangan: ''
  })

  const [details, setDetails] = useState<DetailItem[]>([
    { barang_id: 0, qty: 0, harga: 0 }
  ])

  useEffect(() => {
    fetchPembelians()
    fetchBarangs()
  }, [page])

  const fetchPembelians = async () => {
    setLoading(true)
    try {
      const response = await api.get(`/pembelian?page=${page}&limit=${limit}`)
      setPembelians(response.data.data || [])
      setTotal(response.data.meta?.total || 0)
    } catch (error) {
      console.error('Error fetching pembelian:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchBarangs = async () => {
    try {
      const response = await api.get('/barang?limit=100')
      setBarangs(response.data.data || [])
    } catch (error) {
      console.error('Error fetching barang:', error)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      const payload = {
        ...formData,
        tanggal: formData.tanggal + ':00Z',
        details: details.filter(d => d.barang_id > 0 && d.qty > 0)
      }

      await api.post('/pembelian', payload)
      toast.success('Pembelian created successfully!')
      setShowModal(false)
      resetForm()
      fetchPembelians()
    } catch (error: any) {
      toast.error(error.response?.data?.message || 'Operation failed')
    } finally {
      setLoading(false)
    }
  }

  const addDetail = () => {
    setDetails([...details, { barang_id: 0, qty: 0, harga: 0 }])
  }

  const removeDetail = (index: number) => {
    setDetails(details.filter((_, i) => i !== index))
  }

  const updateDetail = (index: number, field: keyof DetailItem, value: number) => {
    const newDetails = [...details]
    newDetails[index][field] = value
    setDetails(newDetails)
  }

  const resetForm = () => {
    setFormData({
      no_faktur: '',
      tanggal: '',
      supplier: '',
      keterangan: ''
    })
    setDetails([{ barang_id: 0, qty: 0, harga: 0 }])
  }

  const totalPages = Math.ceil(total / limit)

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Pembelian</h1>
        <button
          onClick={() => setShowModal(true)}
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          + Tambah Pembelian
        </button>
      </div>

      {/* Table */}
      <div className="bg-white shadow-md rounded-lg overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">No Faktur</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tanggal</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Supplier</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Total</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Keterangan</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {loading ? (
              <tr>
                <td colSpan={5} className="px-6 py-4 text-center">Loading...</td>
              </tr>
            ) : pembelians.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-6 py-4 text-center">No data found</td>
              </tr>
            ) : (
              pembelians.map((pembelian) => (
                <tr key={pembelian.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{pembelian.no_faktur}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {new Date(pembelian.tanggal).toLocaleDateString('id-ID')}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{pembelian.supplier}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">Rp {pembelian.total.toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm text-gray-500">{pembelian.keterangan}</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      <div className="mt-4 flex justify-between items-center">
        <div className="text-sm text-gray-700">
          Showing {pembelians.length} of {total} results
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => setPage(page - 1)}
            disabled={page === 1}
            className="px-4 py-2 border rounded disabled:opacity-50"
          >
            Previous
          </button>
          <span className="px-4 py-2">Page {page} of {totalPages}</span>
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
          <div className="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75" onClick={() => setShowModal(false)}></div>

            <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-4xl sm:w-full">
              <form onSubmit={handleSubmit}>
                <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                  <h3 className="text-lg font-medium text-gray-900 mb-4">Tambah Pembelian</h3>

                  <div className="grid grid-cols-2 gap-4 mb-4">
                    <div>
                      <label className="block text-gray-700 text-sm font-bold mb-2">No Faktur</label>
                      <input
                        type="text"
                        value={formData.no_faktur}
                        onChange={(e) => setFormData({ ...formData, no_faktur: e.target.value })}
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700"
                        required
                      />
                    </div>
                    <div>
                      <label className="block text-gray-700 text-sm font-bold mb-2">Tanggal</label>
                      <input
                        type="datetime-local"
                        value={formData.tanggal}
                        onChange={(e) => setFormData({ ...formData, tanggal: e.target.value })}
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700"
                        required
                      />
                    </div>
                  </div>

                  <div className="mb-4">
                    <label className="block text-gray-700 text-sm font-bold mb-2">Supplier</label>
                    <input
                      type="text"
                      value={formData.supplier}
                      onChange={(e) => setFormData({ ...formData, supplier: e.target.value })}
                      className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700"
                      required
                    />
                  </div>

                  <div className="mb-4">
                    <label className="block text-gray-700 text-sm font-bold mb-2">Keterangan</label>
                    <input
                      type="text"
                      value={formData.keterangan}
                      onChange={(e) => setFormData({ ...formData, keterangan: e.target.value })}
                      className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700"
                    />
                  </div>

                  <div className="mb-4">
                    <div className="flex justify-between items-center mb-2">
                      <label className="block text-gray-700 text-sm font-bold">Detail Items</label>
                      <button
                        type="button"
                        onClick={addDetail}
                        className="bg-green-500 text-white px-3 py-1 rounded text-sm"
                      >
                        + Add Item
                      </button>
                    </div>

                    {details.map((detail, index) => (
                      <div key={index} className="grid grid-cols-4 gap-2 mb-2">
                        <select
                          value={detail.barang_id}
                          onChange={(e) => updateDetail(index, 'barang_id', Number(e.target.value))}
                          className="border rounded py-2 px-3 text-gray-700"
                          required
                        >
                          <option value={0}>Pilih Barang</option>
                          {barangs.map((barang) => (
                            <option key={barang.id} value={barang.id}>
                              {barang.kode_barang} - {barang.nama_barang}
                            </option>
                          ))}
                        </select>
                        <input
                          type="number"
                          placeholder="Qty"
                          value={detail.qty || ''}
                          onChange={(e) => updateDetail(index, 'qty', Number(e.target.value))}
                          className="border rounded py-2 px-3 text-gray-700"
                          required
                        />
                        <input
                          type="number"
                          placeholder="Harga"
                          value={detail.harga || ''}
                          onChange={(e) => updateDetail(index, 'harga', Number(e.target.value))}
                          className="border rounded py-2 px-3 text-gray-700"
                          required
                        />
                        <button
                          type="button"
                          onClick={() => removeDetail(index)}
                          className="bg-red-500 text-white px-3 py-2 rounded"
                        >
                          Remove
                        </button>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
                  <button
                    type="submit"
                    disabled={loading}
                    className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
                  >
                    {loading ? 'Saving...' : 'Save'}
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      setShowModal(false)
                      resetForm()
                    }}
                    className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
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
