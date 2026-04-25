"use client"

import { useBlockScanner, Transaction, Block } from "@/hooks/useBlockScanner"
import { formatAddress, formatHash, cn } from "@/lib/utils"
import { Activity, Box, Database, Shield, Zap, Search, ChevronRight } from "lucide-react"
import { motion } from "framer-motion"
import Link from "next/link"

export default function Dashboard() {
  const { latestBlock, blocks, transactions, loading } = useBlockScanner()

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="flex flex-col items-center gap-4">
          <div className="w-12 h-12 border-4 border-primary border-t-transparent rounded-full animate-spin" />
          <p className="text-gray-400 font-medium">Scanning FHISH Rollup...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-7xl mx-auto px-6 py-12">
      {/* Header Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-12">
        <StatCard 
          icon={<Box className="text-blue-500" />} 
          label="Latest Block" 
          value={latestBlock.toString()} 
          sub="Mining at ~1s"
        />
        <StatCard 
          icon={<Activity className="text-emerald-500" />} 
          label="Transactions" 
          value={transactions.length.toString()} 
          sub="Recent window"
        />
        <StatCard 
          icon={<Shield className="text-purple-500" />} 
          label="FHE Status" 
          value="Online" 
          sub="Gateway v2.0"
        />
        <StatCard 
          icon={<Zap className="text-amber-500" />} 
          label="Network" 
          value="fhish-1" 
          sub="MiniEVM Stack"
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
        {/* Recent Blocks */}
        <section>
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-bold flex items-center gap-2">
              <Database className="w-5 h-5 text-gray-400" />
              Latest Blocks
            </h2>
          </div>
          <div className="space-y-4">
            {blocks.map((block) => (
              <BlockItem key={block.hash} block={block} />
            ))}
          </div>
        </section>

        {/* Recent Transactions */}
        <section>
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-bold flex items-center gap-2">
              <Zap className="w-5 h-5 text-gray-400" />
              Recent Transactions
            </h2>
          </div>
          <div className="space-y-4">
            {transactions.map((tx) => (
              <TransactionItem key={tx.hash} tx={tx} />
            ))}
          </div>
        </section>
      </div>
    </div>
  )
}

function StatCard({ icon, label, value, sub }: { icon: React.ReactNode, label: string, value: string, subText?: string, sub: string }) {
  return (
    <motion.div 
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="glass p-6 rounded-2xl"
    >
      <div className="flex items-center gap-3 mb-4">
        {icon}
        <span className="text-gray-400 text-sm font-medium">{label}</span>
      </div>
      <div className="text-3xl font-bold mb-1">{value}</div>
      <div className="text-xs text-gray-500">{sub}</div>
    </motion.div>
  )
}

function BlockItem({ block }: { block: Block }) {
  return (
    <div className="glass p-4 rounded-xl flex items-center justify-between hover:border-gray-700 transition-colors group cursor-pointer">
      <div className="flex items-center gap-4">
        <div className="bg-gray-800 p-3 rounded-lg group-hover:bg-blue-900/30 transition-colors">
          <Box className="w-5 h-5 text-gray-400 group-hover:text-blue-400" />
        </div>
        <div>
          <div className="font-bold text-blue-400">#{block.number}</div>
          <div className="text-xs text-gray-500">{new Date(block.timestamp * 1000).toLocaleTimeString()}</div>
        </div>
      </div>
      <div className="text-right">
        <div className="text-sm font-mono text-gray-400">{formatHash(block.hash)}</div>
        <div className="text-xs text-gray-500">{block.transactions.length} txns</div>
      </div>
    </div>
  )
}

function TransactionItem({ tx }: { tx: Transaction }) {
  return (
    <Link href={`/transaction/${tx.hash}`}>
      <div className="glass p-4 rounded-xl flex items-center justify-between hover:border-gray-700 transition-colors group cursor-pointer">
        <div className="flex items-center gap-4">
          <div className={cn(
            "p-3 rounded-lg transition-colors",
            tx.isFHE ? "bg-purple-900/20 text-purple-400" : "bg-gray-800 text-gray-400"
          )}>
            {tx.isFHE ? <Shield className="w-5 h-5" /> : <Search className="w-5 h-5" />}
          </div>
          <div>
            <div className="text-sm font-mono flex items-center gap-2">
              {formatHash(tx.hash)}
              {tx.isFHE && (
                <span className="bg-purple-500/10 text-purple-400 text-[10px] px-1.5 py-0.5 rounded border border-purple-500/20 uppercase tracking-wider font-bold">
                  FHE Call
                </span>
              )}
            </div>
            <div className="text-xs text-gray-500">
              From {formatAddress(tx.from)}
            </div>
          </div>
        </div>
        <div className="text-right">
          <div className="text-sm font-bold text-gray-300">{tx.value} INIT</div>
          <div className="text-xs text-gray-500">Block #{tx.blockNumber}</div>
        </div>
      </div>
    </Link>
  )
}
