import Link from 'next/link'
import { Button } from '@/components/ui/button'

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-b from-gray-50 to-gray-100">
      <div className="text-center space-y-6 p-8">
        <h1 className="text-4xl font-bold text-gray-900">
          FastenMind 緊固件詢報價系統
        </h1>
        <p className="text-xl text-gray-600 max-w-2xl mx-auto">
          專為緊固件產業設計的智能詢價報價管理平台
          <br />
          整合 AI 輔助、N8N 工作流程、全球貿易管理
        </p>
        <div className="flex gap-4 justify-center pt-4">
          <Link href="/login">
            <Button size="lg" className="font-semibold">
              登入系統
            </Button>
          </Link>
          <Link href="/about">
            <Button size="lg" variant="outline" className="font-semibold">
              了解更多
            </Button>
          </Link>
        </div>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mt-16 max-w-6xl mx-auto p-8">
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h3 className="text-xl font-semibold mb-3">🚀 智能詢報價</h3>
          <p className="text-gray-600">
            AI 輔助材料分析、製程成本計算，快速準確完成報價
          </p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h3 className="text-xl font-semibold mb-3">🌍 全球貿易</h3>
          <p className="text-gray-600">
            支援多國關稅計算、Incoterms 管理、海關合規檢查
          </p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h3 className="text-xl font-semibold mb-3">🔄 自動化流程</h3>
          <p className="text-gray-600">
            N8N 整合實現工作流程自動化，提升效率降低錯誤
          </p>
        </div>
      </div>
    </div>
  )
}