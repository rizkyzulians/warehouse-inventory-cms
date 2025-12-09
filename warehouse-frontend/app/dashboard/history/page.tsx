'use client'

import { useEffect, useState } from 'react'
import api from '@/lib/api'

interface HistoryStok {
  id: number
  barang_id: number
  kode_barang: string
  nama_barang: string
  jenis_transaksi: string
  qty: number
  stok_sebelum: number
  stok_sesudah: number
  keterangan: string
  referensi_id: number | null
  referensi_tipe: string
  created_at: string
}

export default function HistoryStokPage() {
  const [histories, setHistories] = useState<HistoryStok[]>([])
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const limit = 10

  useEffect(() => {
    fetchHistories()
  }, [page])

  const fetchHistories = async () => {
    setLoading(true)
    try {
      const response = await api.get(`/stok/history?page=${page}&limit=${limit}`)
      setHistories(response.data.data || [])
      setTotal(response.data.meta?.total || 0)
    } catch (error) {
      console.error('Error fetching stock history:', error)
    } finally {
      setLoading(false)
    }
  }

  const totalPages = Math.ceil(total / limit)

  const getTransactionBadge = (jenis: string) => {
    if (jenis === 'masuk') {
      return (
        <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
          Masuk
        </span>
      )
    } else if (jenis === 'keluar') {
      return (
        <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-red-100 text-red-800">
          Keluar
        </span>
      )
    }
    return <span className="text-gray-500">{jenis}</span>
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Stock History</h1>
      </div>

      {/* Table */}
      <div className="bg-white shadow-md rounded-lg overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tanggal</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Kode Barang</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nama Barang</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Transaksi</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Qty</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Stok Sebelum</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Stok Sesudah</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Keterangan</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {loading ? (
                <tr>
                  <td colSpan={8} className="px-6 py-4 text-center">Loading...</td>
                </tr>
              ) : histories.length === 0 ? (
                <tr>
                  <td colSpan={8} className="px-6 py-4 text-center">No data found</td>
                </tr>
              ) : (
                histories.map((history) => (
                  <tr key={history.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {new Date(history.created_at).toLocaleDateString('id-ID', {
                        year: 'numeric',
                        month: 'short',
                        day: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit'
                      })}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {history.kode_barang}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-900">
                      {history.nama_barang}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      {getTransactionBadge(history.jenis_transaksi)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 font-semibold">
                      {history.jenis_transaksi === 'masuk' ? '+' : '-'}{history.qty}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {history.stok_sebelum}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 font-semibold">
                      {history.stok_sesudah}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500">
                      {history.keterangan}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Pagination */}
      <div className="mt-4 flex justify-between items-center">
        <div className="text-sm text-gray-700">
          Showing {histories.length} of {total} results
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => setPage(page - 1)}
            disabled={page === 1}
            className="px-4 py-2 border rounded disabled:opacity-50 hover:bg-gray-50"
          >
            Previous
          </button>
          <span className="px-4 py-2">
            Page {page} of {totalPages || 1}
          </span>
          <button
            onClick={() => setPage(page + 1)}
            disabled={page === totalPages || totalPages === 0}
            className="px-4 py-2 border rounded disabled:opacity-50 hover:bg-gray-50"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  )
}
