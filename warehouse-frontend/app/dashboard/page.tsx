'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import api from '@/lib/api'

export default function DashboardPage() {
  const [stats, setStats] = useState({
    totalBarang: 0,
    totalStok: 0,
    totalPembelian: 0,
    totalPenjualan: 0
  })

  useEffect(() => {
    fetchStats()
  }, [])

  const fetchStats = async () => {
    try {
      // Fetch from API
      const barangRes = await api.get('/barang?limit=1')
      const barangStokRes = await api.get('/barang/stok?limit=100')
      const pembelianRes = await api.get('/pembelian?limit=1')
      const penjualanRes = await api.get('/penjualan?limit=1')

      // Count items that have stock (qty_akhir > 0)
      const itemsInStock = barangStokRes.data.data?.filter((item: any) => item.qty_akhir > 0).length || 0

      setStats({
        totalBarang: barangRes.data.meta?.total || 0,
        totalStok: itemsInStock,
        totalPembelian: pembelianRes.data.meta?.total || 0,
        totalPenjualan: penjualanRes.data.meta?.total || 0
      })
    } catch (error) {
      console.error('Error fetching stats:', error)
    }
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Dashboard</h1>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="rounded-md bg-blue-500 p-3">
                  <svg className="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                  </svg>
                </div>
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">Total Barang</dt>
                  <dd className="text-3xl font-semibold text-gray-900">{stats.totalBarang}</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="rounded-md bg-green-500 p-3">
                  <svg className="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                  </svg>
                </div>
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">Items in Stock</dt>
                  <dd className="text-3xl font-semibold text-gray-900">{stats.totalStok}</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="rounded-md bg-yellow-500 p-3">
                  <svg className="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
                  </svg>
                </div>
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">Pembelian</dt>
                  <dd className="text-3xl font-semibold text-gray-900">{stats.totalPembelian}</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="rounded-md bg-purple-500 p-3">
                  <svg className="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">Penjualan</dt>
                  <dd className="text-3xl font-semibold text-gray-900">{stats.totalPenjualan}</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Quick Links */}
      <div className="bg-white shadow rounded-lg p-6">
        <h2 className="text-xl font-bold text-gray-900 mb-4">Quick Actions</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <Link href="/dashboard/barang" className="text-center p-4 bg-blue-50 hover:bg-blue-100 rounded-lg transition">
            <div className="text-blue-600 font-semibold">Master Barang</div>
          </Link>
          <Link href="/dashboard/stok" className="text-center p-4 bg-green-50 hover:bg-green-100 rounded-lg transition">
            <div className="text-green-600 font-semibold">Stock</div>
          </Link>
          <Link href="/dashboard/pembelian" className="text-center p-4 bg-yellow-50 hover:bg-yellow-100 rounded-lg transition">
            <div className="text-yellow-600 font-semibold">Pembelian</div>
          </Link>
          <Link href="/dashboard/penjualan" className="text-center p-4 bg-purple-50 hover:bg-purple-100 rounded-lg transition">
            <div className="text-purple-600 font-semibold">Penjualan</div>
          </Link>
        </div>
      </div>
    </div>
  )
}
