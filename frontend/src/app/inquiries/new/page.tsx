'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { format } from 'date-fns'
import DashboardLayout from '@/components/layout/DashboardLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useToast } from '@/components/ui/use-toast'
import { Upload, X, Loader2 } from 'lucide-react'
import inquiryService from '@/services/inquiry.service'

const inquirySchema = z.object({
  customer_id: z.string().min(1, '請選擇客戶'),
  product_category: z.string().min(1, '請選擇產品分類'),
  product_name: z.string().min(1, '請輸入產品名稱'),
  quantity: z.number().min(1, '數量必須大於 0'),
  unit: z.string().min(1, '請選擇單位'),
  required_date: z.string().min(1, '請選擇交期'),
  incoterm: z.string().min(1, '請選擇交易條件'),
  destination_port: z.string().optional(),
  destination_address: z.string().optional(),
  payment_terms: z.string().optional(),
  special_requirements: z.string().optional(),
})

type InquiryFormData = z.infer<typeof inquirySchema>

export default function NewInquiryPage() {
  const router = useRouter()
  const { toast } = useToast()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [uploadedFiles, setUploadedFiles] = useState<string[]>([])

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    watch,
  } = useForm<InquiryFormData>({
    resolver: zodResolver(inquirySchema),
    defaultValues: {
      quantity: 1,
      unit: 'pcs',
    },
  })

  const onSubmit = async (data: InquiryFormData) => {
    if (uploadedFiles.length === 0) {
      toast({
        title: '請上傳圖紙',
        description: '至少需要上傳一個圖紙檔案',
        variant: 'destructive',
      })
      return
    }

    setIsSubmitting(true)
    try {
      await inquiryService.create({
        ...data,
        drawing_files: uploadedFiles,
      })

      toast({
        title: '詢價單建立成功',
        description: '詢價單已成功建立並發送給工程部門',
      })

      router.push('/inquiries')
    } catch (error: any) {
      toast({
        title: '建立失敗',
        description: error.response?.data?.message || '建立詢價單時發生錯誤',
        variant: 'destructive',
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files
    if (!files) return

    for (const file of Array.from(files)) {
      try {
        const result = await inquiryService.uploadDrawing(file)
        setUploadedFiles((prev) => [...prev, result.url])
        toast({
          title: '檔案上傳成功',
          description: `${file.name} 已上傳`,
        })
      } catch (error) {
        toast({
          title: '檔案上傳失敗',
          description: `${file.name} 上傳失敗`,
          variant: 'destructive',
        })
      }
    }
  }

  const removeFile = (index: number) => {
    setUploadedFiles((prev) => prev.filter((_, i) => i !== index))
  }

  return (
    <DashboardLayout>
      <div className="max-w-4xl mx-auto space-y-6">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-bold text-gray-900">新增詢價單</h1>
          <p className="mt-2 text-gray-600">填寫詢價資訊並上傳相關圖紙</p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          {/* Customer Information */}
          <Card>
            <CardHeader>
              <CardTitle>客戶資訊</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="customer_id">客戶 *</Label>
                  <Select onValueChange={(value) => setValue('customer_id', value)}>
                    <SelectTrigger>
                      <SelectValue placeholder="選擇客戶" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="cust-001">台灣機械股份有限公司</SelectItem>
                      <SelectItem value="cust-002">Global Auto Parts GmbH</SelectItem>
                      <SelectItem value="cust-003">American Tools Inc.</SelectItem>
                    </SelectContent>
                  </Select>
                  {errors.customer_id && (
                    <p className="text-sm text-red-500">{errors.customer_id.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="required_date">要求交期 *</Label>
                  <Input
                    type="date"
                    {...register('required_date')}
                    min={format(new Date(), 'yyyy-MM-dd')}
                  />
                  {errors.required_date && (
                    <p className="text-sm text-red-500">{errors.required_date.message}</p>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Product Information */}
          <Card>
            <CardHeader>
              <CardTitle>產品資訊</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="product_category">產品分類 *</Label>
                  <Select onValueChange={(value) => setValue('product_category', value)}>
                    <SelectTrigger>
                      <SelectValue placeholder="選擇產品分類" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="螺絲">螺絲</SelectItem>
                      <SelectItem value="螺帽">螺帽</SelectItem>
                      <SelectItem value="華司">華司</SelectItem>
                      <SelectItem value="特殊件">特殊件</SelectItem>
                    </SelectContent>
                  </Select>
                  {errors.product_category && (
                    <p className="text-sm text-red-500">{errors.product_category.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="product_name">產品名稱 *</Label>
                  <Input
                    {...register('product_name')}
                    placeholder="例如：M8x30 六角螺栓"
                  />
                  {errors.product_name && (
                    <p className="text-sm text-red-500">{errors.product_name.message}</p>
                  )}
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="quantity">數量 *</Label>
                  <Input
                    type="number"
                    {...register('quantity', { valueAsNumber: true })}
                    min="1"
                  />
                  {errors.quantity && (
                    <p className="text-sm text-red-500">{errors.quantity.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="unit">單位 *</Label>
                  <Select 
                    defaultValue="pcs"
                    onValueChange={(value) => setValue('unit', value)}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="pcs">件 (pcs)</SelectItem>
                      <SelectItem value="kg">公斤 (kg)</SelectItem>
                      <SelectItem value="set">套 (set)</SelectItem>
                    </SelectContent>
                  </Select>
                  {errors.unit && (
                    <p className="text-sm text-red-500">{errors.unit.message}</p>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Drawing Files */}
          <Card>
            <CardHeader>
              <CardTitle>圖紙檔案</CardTitle>
              <CardDescription>
                請上傳產品圖紙或規格書 (支援 PDF, DWG, JPG, PNG 格式)
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center">
                  <Upload className="mx-auto h-12 w-12 text-gray-400" />
                  <div className="mt-4">
                    <label htmlFor="file-upload" className="cursor-pointer">
                      <span className="text-sm font-medium text-blue-600 hover:text-blue-500">
                        點擊上傳檔案
                      </span>
                      <input
                        id="file-upload"
                        type="file"
                        className="sr-only"
                        multiple
                        accept=".pdf,.dwg,.jpg,.jpeg,.png"
                        onChange={handleFileUpload}
                      />
                    </label>
                    <p className="text-xs text-gray-500 mt-1">
                      或拖放檔案到此處
                    </p>
                  </div>
                </div>

                {uploadedFiles.length > 0 && (
                  <div className="space-y-2">
                    {uploadedFiles.map((file, index) => (
                      <div
                        key={index}
                        className="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
                      >
                        <span className="text-sm text-gray-700">
                          檔案 {index + 1}
                        </span>
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => removeFile(index)}
                        >
                          <X className="h-4 w-4" />
                        </Button>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </CardContent>
          </Card>

          {/* Trade Terms */}
          <Card>
            <CardHeader>
              <CardTitle>交易條件</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="incoterm">國際貿易條件 *</Label>
                  <Select onValueChange={(value) => setValue('incoterm', value)}>
                    <SelectTrigger>
                      <SelectValue placeholder="選擇 Incoterm" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="EXW">EXW - 工廠交貨</SelectItem>
                      <SelectItem value="FOB">FOB - 裝運港船上交貨</SelectItem>
                      <SelectItem value="CIF">CIF - 成本保險費加運費</SelectItem>
                      <SelectItem value="DDP">DDP - 完稅後交貨</SelectItem>
                    </SelectContent>
                  </Select>
                  {errors.incoterm && (
                    <p className="text-sm text-red-500">{errors.incoterm.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="destination_port">目的港/地址</Label>
                  <Input
                    {...register('destination_port')}
                    placeholder="例如：Hamburg Port"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="payment_terms">付款條件</Label>
                <Input
                  {...register('payment_terms')}
                  placeholder="例如：T/T 30 days"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="special_requirements">特殊要求</Label>
                <Textarea
                  {...register('special_requirements')}
                  placeholder="請輸入任何特殊要求或備註..."
                  rows={4}
                />
              </div>
            </CardContent>
          </Card>

          {/* Actions */}
          <div className="flex justify-end gap-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => router.push('/inquiries')}
            >
              取消
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  提交中...
                </>
              ) : (
                '提交詢價單'
              )}
            </Button>
          </div>
        </form>
      </div>
    </DashboardLayout>
  )
}