'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import api from '@/lib/api'

interface StokBarang {
  id: number
  barang_id: number
  kode_barang: string
  nama_barang: string
  kategori: string
  satuan: string
  harga_beli: number
  harga_jual: number
  qty_masuk: number
  qty_keluar: number
  qty_akhir: number
}

export default function StokPage() {
  const router = useRouter()
  const [stoks, setStoks] = useState<StokBarang[]>([])
  const [loading, setLoading] = useState(false)
  const [search, setSearch] = useState('')
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const limit = 10

  useEffect(() => {
    fetchStoks()
  }, [page, search])

  const fetchStoks = async () => {
    setLoading(true)
    try {
      const response = await api.get(`/barang/stok?search=${search}&page=${page}&limit=${limit}`)
      console.log('Stok Response:', response.data)
      setStoks(response.data.data || [])
      setTotal(response.data.meta?.total || 0)
    } catch (error) {
      console.error('Error fetching stok:', error)
    } finally {
      setLoading(false)
    }
  }

  const totalPages = Math.ceil(total / limit)

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Laporan Stock</h1>
        <div className="flex gap-2">
          <button
            onClick={() => router.push('/dashboard/history')}
            className="bg-indigo-500 hover:bg-indigo-700 text-white font-bold py-2 px-4 rounded"
          >
            View History
          </button>
          <button
            onClick={fetchStoks}
            className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
          >
            Refresh
          </button>
        </div>
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
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Qty Masuk</th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Qty Keluar</th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Stock Akhir</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {loading ? (
              <tr>
                <td colSpan={7} className="px-6 py-4 text-center">Loading...</td>
              </tr>
            ) : stoks.length === 0 ? (
              <tr>
                <td colSpan={7} className="px-6 py-4 text-center">No data found</td>
              </tr>
            ) : (
              stoks.map((stok) => (
                <tr key={stok.barang_id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{stok.kode_barang}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{stok.nama_barang}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{stok.kategori}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{stok.satuan}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 text-right">{stok.qty_masuk}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 text-right">{stok.qty_keluar}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-bold text-gray-900 text-right">
                    <span className={stok.qty_akhir < 10 ? 'text-red-600' : 'text-green-600'}>
                      {stok.qty_akhir}
                    </span>
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
          Showing {stoks.length} of {total} results
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
    </div>
  )
}
