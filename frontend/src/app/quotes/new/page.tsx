'use client';

import { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useQuery, useMutation } from '@tanstack/react-query';
import { Plus, Trash2, Save, ArrowLeft } from 'lucide-react';
import { quoteService } from '@/services/quote.service';
import { customerService } from '@/services/customer.service';
import { inquiryService } from '@/services/inquiry.service';
import { CreateQuoteRequest, QuoteItemRequest, QuoteTermRequest } from '@/types/quote';
import { useAuth } from '@/contexts/AuthContext';
import { toast } from 'react-hot-toast';
import LoadingSpinner from '@/components/LoadingSpinner';
import DashboardLayout from '@/components/layout/DashboardLayout';

export default function NewQuotePage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { user } = useAuth();
  const inquiryId = searchParams.get('inquiry_id');

  const [formData, setFormData] = useState<CreateQuoteRequest>({
    inquiry_id: '',
    customer_id: '',
    validity_days: 30,
    payment_terms: 'T/T 30 days',
    delivery_terms: 'FOB',
    remarks: '',
    items: [],
    terms: [],
    use_template: true,
  });

  const [currentItem, setCurrentItem] = useState<QuoteItemRequest>({
    product_name: '',
    product_specs: '',
    quantity: 1,
    unit: 'PCS',
    unit_price: 0,
    notes: ''
  });

  const [currentTerm, setCurrentTerm] = useState<QuoteTermRequest>({
    term_type: '付款條件',
    term_content: '',
    sort_order: 1
  });

  // Fetch inquiry details if inquiry_id is provided
  const { data: inquiry } = useQuery({
    queryKey: ['inquiry', inquiryId],
    queryFn: () => inquiryService.getInquiry(inquiryId!),
    enabled: !!inquiryId,
  });

  // Fetch customers
  const { data: customers } = useQuery({
    queryKey: ['customers'],
    queryFn: () => customerService.getCustomers({ page_size: 100 }),
  });

  // Fetch quote templates
  const { data: templates } = useQuery({
    queryKey: ['quote-templates'],
    queryFn: () => quoteService.getQuoteTemplates(),
  });

  useEffect(() => {
    if (inquiry) {
      setFormData(prev => ({
        ...prev,
        inquiry_id: inquiry.id,
        customer_id: inquiry.customer_id,
      }));

      // Pre-fill first item with inquiry details
      if (formData.items.length === 0) {
        setCurrentItem({
          product_name: inquiry.product_name || '',
          product_specs: inquiry.product_specs || '',
          quantity: inquiry.quantity || 1,
          unit: inquiry.unit || 'PCS',
          unit_price: 0,
          notes: ''
        });
      }
    }
  }, [inquiry]);

  const createQuoteMutation = useMutation({
    mutationFn: (data: CreateQuoteRequest) => quoteService.createQuote(data),
    onSuccess: (quote) => {
      toast.success('報價單建立成功');
      router.push(`/quotes/${quote.id}`);
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || '建立報價單失敗');
    },
  });

  const handleAddItem = () => {
    if (!currentItem.product_name || currentItem.quantity <= 0 || currentItem.unit_price <= 0) {
      toast.error('請填寫完整的產品資訊');
      return;
    }

    setFormData(prev => ({
      ...prev,
      items: [...prev.items, currentItem]
    }));

    // Reset current item
    setCurrentItem({
      product_name: '',
      product_specs: '',
      quantity: 1,
      unit: 'PCS',
      unit_price: 0,
      notes: ''
    });
  };

  const handleRemoveItem = (index: number) => {
    setFormData(prev => ({
      ...prev,
      items: prev.items.filter((_, i) => i !== index)
    }));
  };

  const handleAddTerm = () => {
    if (!currentTerm.term_content) {
      toast.error('請填寫條款內容');
      return;
    }

    setFormData(prev => ({
      ...prev,
      terms: [...prev.terms, { ...currentTerm, sort_order: prev.terms.length + 1 }]
    }));

    // Reset current term
    setCurrentTerm({
      term_type: '付款條件',
      term_content: '',
      sort_order: 1
    });
  };

  const handleRemoveTerm = (index: number) => {
    setFormData(prev => ({
      ...prev,
      terms: prev.terms.filter((_, i) => i !== index)
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.customer_id) {
      toast.error('請選擇客戶');
      return;
    }

    if (formData.items.length === 0) {
      toast.error('請至少添加一個報價項目');
      return;
    }

    createQuoteMutation.mutate(formData);
  };

  return (
    <DashboardLayout>
      <div className="p-6">
        <div className="mb-6">
          <div className="flex items-center gap-4">
            <button
              onClick={() => router.back()}
              className="p-2 hover:bg-gray-100 rounded-md"
            >
              <ArrowLeft className="h-5 w-5" />
            </button>
            <div>
              <h1 className="text-2xl font-bold text-gray-900">新增報價單</h1>
              <p className="text-gray-600 mt-1">填寫報價單資訊並添加報價項目</p>
            </div>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* 基本資訊 */}
          <div className="bg-white shadow rounded-lg p-6">
            <h2 className="text-lg font-semibold mb-4">基本資訊</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {inquiryId && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    詢價單號
                  </label>
                  <input
                    type="text"
                    value={inquiry?.inquiry_no || ''}
                    disabled
                    className="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50"
                  />
                </div>
              )}

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  客戶 <span className="text-red-500">*</span>
                </label>
                <select
                  value={formData.customer_id}
                  onChange={(e) => setFormData(prev => ({ ...prev, customer_id: e.target.value }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                >
                  <option value="">選擇客戶</option>
                  {customers?.data.map((customer) => (
                    <option key={customer.id} value={customer.id}>
                      {customer.name} ({customer.customer_code})
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  有效天數
                </label>
                <input
                  type="number"
                  value={formData.validity_days}
                  onChange={(e) => setFormData(prev => ({ ...prev, validity_days: parseInt(e.target.value) }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  min={1}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  付款條件
                </label>
                <input
                  type="text"
                  value={formData.payment_terms}
                  onChange={(e) => setFormData(prev => ({ ...prev, payment_terms: e.target.value }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  交貨條件
                </label>
                <input
                  type="text"
                  value={formData.delivery_terms}
                  onChange={(e) => setFormData(prev => ({ ...prev, delivery_terms: e.target.value }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  備註
                </label>
                <textarea
                  value={formData.remarks}
                  onChange={(e) => setFormData(prev => ({ ...prev, remarks: e.target.value }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  rows={3}
                />
              </div>
            </div>
          </div>

          {/* 報價項目 */}
          <div className="bg-white shadow rounded-lg p-6">
            <h2 className="text-lg font-semibold mb-4">報價項目</h2>
            
            {/* 新增項目表單 */}
            <div className="bg-gray-50 p-4 rounded-md mb-4">
              <div className="grid grid-cols-1 md:grid-cols-6 gap-3">
                <div className="md:col-span-2">
                  <input
                    type="text"
                    placeholder="產品名稱"
                    value={currentItem.product_name}
                    onChange={(e) => setCurrentItem(prev => ({ ...prev, product_name: e.target.value }))}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  />
                </div>
                <div className="md:col-span-2">
                  <input
                    type="text"
                    placeholder="規格"
                    value={currentItem.product_specs}
                    onChange={(e) => setCurrentItem(prev => ({ ...prev, product_specs: e.target.value }))}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  />
                </div>
                <div>
                  <input
                    type="number"
                    placeholder="數量"
                    value={currentItem.quantity}
                    onChange={(e) => setCurrentItem(prev => ({ ...prev, quantity: parseInt(e.target.value) }))}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                    min={1}
                  />
                </div>
                <div>
                  <input
                    type="text"
                    placeholder="單位"
                    value={currentItem.unit}
                    onChange={(e) => setCurrentItem(prev => ({ ...prev, unit: e.target.value }))}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  />
                </div>
                <div>
                  <input
                    type="number"
                    placeholder="單價"
                    value={currentItem.unit_price}
                    onChange={(e) => setCurrentItem(prev => ({ ...prev, unit_price: parseFloat(e.target.value) }))}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                    step="0.0001"
                    min={0}
                  />
                </div>
                <div className="md:col-span-5">
                  <input
                    type="text"
                    placeholder="備註"
                    value={currentItem.notes}
                    onChange={(e) => setCurrentItem(prev => ({ ...prev, notes: e.target.value }))}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  />
                </div>
                <div>
                  <button
                    type="button"
                    onClick={handleAddItem}
                    className="w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 flex items-center justify-center gap-2"
                  >
                    <Plus className="h-4 w-4" />
                    新增
                  </button>
                </div>
              </div>
            </div>

            {/* 項目列表 */}
            {formData.items.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        產品名稱
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        規格
                      </th>
                      <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                        數量
                      </th>
                      <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                        單位
                      </th>
                      <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                        單價
                      </th>
                      <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                        總價
                      </th>
                      <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                        操作
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {formData.items.map((item, index) => (
                      <tr key={index}>
                        <td className="px-4 py-3 text-sm text-gray-900">
                          {item.product_name}
                        </td>
                        <td className="px-4 py-3 text-sm text-gray-500">
                          {item.product_specs || '-'}
                        </td>
                        <td className="px-4 py-3 text-sm text-gray-900 text-right">
                          {item.quantity.toLocaleString()}
                        </td>
                        <td className="px-4 py-3 text-sm text-gray-900 text-center">
                          {item.unit}
                        </td>
                        <td className="px-4 py-3 text-sm text-gray-900 text-right">
                          ${item.unit_price.toFixed(4)}
                        </td>
                        <td className="px-4 py-3 text-sm font-medium text-gray-900 text-right">
                          ${(item.quantity * item.unit_price).toFixed(2)}
                        </td>
                        <td className="px-4 py-3 text-center">
                          <button
                            type="button"
                            onClick={() => handleRemoveItem(index)}
                            className="text-red-600 hover:text-red-900"
                          >
                            <Trash2 className="h-4 w-4" />
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                  <tfoot className="bg-gray-50">
                    <tr>
                      <td colSpan={5} className="px-4 py-3 text-right text-sm font-medium text-gray-900">
                        總計
                      </td>
                      <td className="px-4 py-3 text-right text-sm font-bold text-gray-900">
                        ${formData.items.reduce((sum, item) => sum + item.quantity * item.unit_price, 0).toFixed(2)}
                      </td>
                      <td></td>
                    </tr>
                  </tfoot>
                </table>
              </div>
            ) : (
              <div className="text-center py-8 text-gray-500">
                尚未添加報價項目
              </div>
            )}
          </div>

          {/* 條款設定 */}
          <div className="bg-white shadow rounded-lg p-6">
            <h2 className="text-lg font-semibold mb-4">條款設定</h2>
            
            <div className="mb-4">
              <label className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={formData.use_template}
                  onChange={(e) => setFormData(prev => ({ ...prev, use_template: e.target.checked }))}
                  className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                />
                <span className="text-sm font-medium text-gray-700">使用預設條款模板</span>
              </label>
            </div>

            {!formData.use_template && (
              <>
                {/* 新增條款表單 */}
                <div className="bg-gray-50 p-4 rounded-md mb-4">
                  <div className="grid grid-cols-1 md:grid-cols-4 gap-3">
                    <div>
                      <select
                        value={currentTerm.term_type}
                        onChange={(e) => setCurrentTerm(prev => ({ ...prev, term_type: e.target.value }))}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md"
                      >
                        <option value="付款條件">付款條件</option>
                        <option value="交貨條件">交貨條件</option>
                        <option value="品質要求">品質要求</option>
                        <option value="其他條款">其他條款</option>
                      </select>
                    </div>
                    <div className="md:col-span-2">
                      <textarea
                        placeholder="條款內容"
                        value={currentTerm.term_content}
                        onChange={(e) => setCurrentTerm(prev => ({ ...prev, term_content: e.target.value }))}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md"
                        rows={2}
                      />
                    </div>
                    <div>
                      <button
                        type="button"
                        onClick={handleAddTerm}
                        className="w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 flex items-center justify-center gap-2 h-full"
                      >
                        <Plus className="h-4 w-4" />
                        新增
                      </button>
                    </div>
                  </div>
                </div>

                {/* 條款列表 */}
                {formData.terms.length > 0 && (
                  <div className="space-y-3">
                    {formData.terms.map((term, index) => (
                      <div key={index} className="border rounded-lg p-4 flex justify-between items-start">
                        <div className="flex-1">
                          <h5 className="font-medium text-sm text-gray-700 mb-1">{term.term_type}</h5>
                          <p className="text-sm text-gray-600 whitespace-pre-wrap">{term.term_content}</p>
                        </div>
                        <button
                          type="button"
                          onClick={() => handleRemoveTerm(index)}
                          className="ml-4 text-red-600 hover:text-red-900"
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    ))}
                  </div>
                )}
              </>
            )}
          </div>

          {/* 操作按鈕 */}
          <div className="flex justify-end gap-4">
            <button
              type="button"
              onClick={() => router.back()}
              className="px-6 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
            >
              取消
            </button>
            <button
              type="submit"
              disabled={createQuoteMutation.isLoading}
              className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 flex items-center gap-2"
            >
              {createQuoteMutation.isLoading ? (
                <>
                  <LoadingSpinner className="h-4 w-4" />
                  建立中...
                </>
              ) : (
                <>
                  <Save className="h-4 w-4" />
                  建立報價單
                </>
              )}
            </button>
          </div>
        </form>
      </div>
    </DashboardLayout>
  );
}