'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog';
import { toast } from '@/components/ui/use-toast';
import { 
  Calculator, 
  Package, 
  Settings, 
  FileSpreadsheet,
  DollarSign,
  TrendingUp,
  Clock,
  Layers
} from 'lucide-react';

interface ProcessCostTemplate {
  id: string;
  name: string;
  description: string;
  process_type: string;
  category: string;
  total_cost: number;
  is_active: boolean;
}

interface MaterialCost {
  id: string;
  name: string;
  type: string;
  specification: string;
  unit_price: number;
  currency: string;
  unit: string;
  supplier: string;
  price_change_reason?: string;
}

interface ProcessingRate {
  id: string;
  process_type: string;
  equipment_name: string;
  hourly_rate: number;
  setup_cost: number;
  minimum_charge: number;
  currency: string;
}

export default function ProcessCostsPage() {
  const [activeTab, setActiveTab] = useState('calculate');
  const [templates, setTemplates] = useState<ProcessCostTemplate[]>([]);
  const [materials, setMaterials] = useState<MaterialCost[]>([]);
  const [processingRates, setProcessingRates] = useState<ProcessingRate[]>([]);
  const [loading, setLoading] = useState(false);
  
  // 成本計算表單
  const [costForm, setCostForm] = useState({
    product_spec: {
      length: 0,
      width: 0,
      height: 0,
      diameter: 0,
      thickness: 0,
      weight: 0,
      complexity: 'medium'
    },
    material_id: '',
    material_utilization: 90,
    quantity: 1,
    processes: [] as any[],
    surface_treatment: '',
    overhead_rate: 15,
    profit_margin: 20,
    base_currency: 'TWD',
    target_currency: 'TWD'
  });
  
  // 計算結果
  const [calculationResult, setCalculationResult] = useState<any>(null);
  
  // 對話框
  const [showMaterialDialog, setShowMaterialDialog] = useState(false);
  const [showProcessDialog, setShowProcessDialog] = useState(false);
  const [selectedMaterial, setSelectedMaterial] = useState<MaterialCost | null>(null);
  
  useEffect(() => {
    fetchData();
  }, [activeTab]);

  const fetchData = async () => {
    setLoading(true);
    try {
      if (activeTab === 'templates') {
        await fetchTemplates();
      } else if (activeTab === 'materials') {
        await fetchMaterials();
      } else if (activeTab === 'rates') {
        await fetchProcessingRates();
      }
    } catch (error) {
      console.error('Failed to fetch data:', error);
      toast({
        title: '錯誤',
        description: '載入資料失敗',
        variant: 'destructive'
      });
    } finally {
      setLoading(false);
    }
  };

  const fetchTemplates = async () => {
    const response = await fetch('/api/v1/process-costs/templates');
    if (response.ok) {
      const data = await response.json();
      setTemplates(data.data || []);
    }
  };

  const fetchMaterials = async () => {
    const response = await fetch('/api/v1/process-costs/materials');
    if (response.ok) {
      const data = await response.json();
      setMaterials(data.data || []);
    }
  };

  const fetchProcessingRates = async () => {
    const response = await fetch('/api/v1/process-costs/processing-rates');
    if (response.ok) {
      const data = await response.json();
      setProcessingRates(data.data || []);
    }
  };

  const handleCalculate = async () => {
    try {
      const response = await fetch('/api/v1/process-costs/calculate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(costForm)
      });

      if (response.ok) {
        const data = await response.json();
        setCalculationResult(data.data);
        toast({
          title: '成功',
          description: '成本計算完成'
        });
      } else {
        throw new Error('計算失敗');
      }
    } catch (error) {
      toast({
        title: '錯誤',
        description: '成本計算失敗',
        variant: 'destructive'
      });
    }
  };

  const handleAddProcess = () => {
    setCostForm({
      ...costForm,
      processes: [
        ...costForm.processes,
        {
          process_type: '',
          equipment_id: '',
          parameters: {},
          sequence: costForm.processes.length + 1
        }
      ]
    });
  };

  const handleUpdateMaterial = async () => {
    if (!selectedMaterial) return;

    try {
      const response = await fetch(`/api/v1/process-costs/materials/${selectedMaterial.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(selectedMaterial)
      });

      if (response.ok) {
        toast({
          title: '成功',
          description: '材料成本更新成功'
        });
        setShowMaterialDialog(false);
        fetchMaterials();
      } else {
        throw new Error('更新失敗');
      }
    } catch (error) {
      toast({
        title: '錯誤',
        description: '更新材料成本失敗',
        variant: 'destructive'
      });
    }
  };

  const formatCurrency = (amount: number, currency: string = 'TWD') => {
    return new Intl.NumberFormat('zh-TW', {
      style: 'currency',
      currency: currency
    }).format(amount);
  };

  return (
    <div className="container mx-auto p-6 space-y-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">製程成本計算</h1>
        <Button onClick={() => setActiveTab('settings')} variant="outline">
          <Settings className="mr-2 h-4 w-4" />
          成本設定
        </Button>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="calculate">成本計算</TabsTrigger>
          <TabsTrigger value="templates">成本模板</TabsTrigger>
          <TabsTrigger value="materials">材料成本</TabsTrigger>
          <TabsTrigger value="rates">加工費率</TabsTrigger>
        </TabsList>

        <TabsContent value="calculate" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>產品規格</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label>長度 (mm)</Label>
                    <Input
                      type="number"
                      value={costForm.product_spec.length}
                      onChange={(e) => setCostForm({
                        ...costForm,
                        product_spec: {
                          ...costForm.product_spec,
                          length: parseFloat(e.target.value) || 0
                        }
                      })}
                    />
                  </div>
                  <div>
                    <Label>寬度 (mm)</Label>
                    <Input
                      type="number"
                      value={costForm.product_spec.width}
                      onChange={(e) => setCostForm({
                        ...costForm,
                        product_spec: {
                          ...costForm.product_spec,
                          width: parseFloat(e.target.value) || 0
                        }
                      })}
                    />
                  </div>
                  <div>
                    <Label>高度 (mm)</Label>
                    <Input
                      type="number"
                      value={costForm.product_spec.height}
                      onChange={(e) => setCostForm({
                        ...costForm,
                        product_spec: {
                          ...costForm.product_spec,
                          height: parseFloat(e.target.value) || 0
                        }
                      })}
                    />
                  </div>
                  <div>
                    <Label>直徑 (mm)</Label>
                    <Input
                      type="number"
                      value={costForm.product_spec.diameter}
                      onChange={(e) => setCostForm({
                        ...costForm,
                        product_spec: {
                          ...costForm.product_spec,
                          diameter: parseFloat(e.target.value) || 0
                        }
                      })}
                    />
                  </div>
                </div>
                
                <div>
                  <Label>複雜度</Label>
                  <Select
                    value={costForm.product_spec.complexity}
                    onValueChange={(value) => setCostForm({
                      ...costForm,
                      product_spec: {
                        ...costForm.product_spec,
                        complexity: value
                      }
                    })}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="low">低</SelectItem>
                      <SelectItem value="medium">中</SelectItem>
                      <SelectItem value="high">高</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>材料與數量</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label>材料</Label>
                  <Select
                    value={costForm.material_id}
                    onValueChange={(value) => setCostForm({...costForm, material_id: value})}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="選擇材料" />
                    </SelectTrigger>
                    <SelectContent>
                      {materials.map((material) => (
                        <SelectItem key={material.id} value={material.id}>
                          {material.name} - {material.specification}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                
                <div>
                  <Label>材料利用率 (%)</Label>
                  <Input
                    type="number"
                    value={costForm.material_utilization}
                    onChange={(e) => setCostForm({
                      ...costForm,
                      material_utilization: parseFloat(e.target.value) || 90
                    })}
                  />
                </div>
                
                <div>
                  <Label>數量</Label>
                  <Input
                    type="number"
                    value={costForm.quantity}
                    onChange={(e) => setCostForm({
                      ...costForm,
                      quantity: parseInt(e.target.value) || 1
                    })}
                  />
                </div>
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>加工製程</CardTitle>
              <Button onClick={handleAddProcess} size="sm">
                新增製程
              </Button>
            </CardHeader>
            <CardContent>
              {costForm.processes.length === 0 ? (
                <p className="text-muted-foreground text-center py-4">
                  尚未添加製程步驟
                </p>
              ) : (
                <div className="space-y-2">
                  {costForm.processes.map((process, index) => (
                    <div key={index} className="flex gap-2 items-center">
                      <span className="text-sm font-medium">步驟 {index + 1}</span>
                      <Select
                        value={process.process_type}
                        onValueChange={(value) => {
                          const updatedProcesses = [...costForm.processes];
                          updatedProcesses[index].process_type = value;
                          setCostForm({...costForm, processes: updatedProcesses});
                        }}
                      >
                        <SelectTrigger className="flex-1">
                          <SelectValue placeholder="選擇製程" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="turning">車削</SelectItem>
                          <SelectItem value="milling">銑削</SelectItem>
                          <SelectItem value="drilling">鑽孔</SelectItem>
                          <SelectItem value="grinding">研磨</SelectItem>
                        </SelectContent>
                      </Select>
                      <Button
                        size="sm"
                        variant="destructive"
                        onClick={() => {
                          const updatedProcesses = costForm.processes.filter((_, i) => i !== index);
                          setCostForm({...costForm, processes: updatedProcesses});
                        }}
                      >
                        刪除
                      </Button>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>成本參數</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>管理費率 (%)</Label>
                  <Input
                    type="number"
                    value={costForm.overhead_rate}
                    onChange={(e) => setCostForm({
                      ...costForm,
                      overhead_rate: parseFloat(e.target.value) || 15
                    })}
                  />
                </div>
                <div>
                  <Label>利潤率 (%)</Label>
                  <Input
                    type="number"
                    value={costForm.profit_margin}
                    onChange={(e) => setCostForm({
                      ...costForm,
                      profit_margin: parseFloat(e.target.value) || 20
                    })}
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          <div className="flex justify-end">
            <Button onClick={handleCalculate} size="lg">
              <Calculator className="mr-2 h-4 w-4" />
              計算成本
            </Button>
          </div>

          {calculationResult && (
            <Card>
              <CardHeader>
                <CardTitle>計算結果</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                  <div className="space-y-2">
                    <p className="text-sm text-muted-foreground">材料成本</p>
                    <p className="text-2xl font-bold">
                      {formatCurrency(calculationResult.material_cost)}
                    </p>
                  </div>
                  <div className="space-y-2">
                    <p className="text-sm text-muted-foreground">加工成本</p>
                    <p className="text-2xl font-bold">
                      {formatCurrency(calculationResult.processing_cost)}
                    </p>
                  </div>
                  <div className="space-y-2">
                    <p className="text-sm text-muted-foreground">管理費用</p>
                    <p className="text-2xl font-bold">
                      {formatCurrency(calculationResult.overhead_cost)}
                    </p>
                  </div>
                  <div className="space-y-2">
                    <p className="text-sm text-muted-foreground">總成本</p>
                    <p className="text-2xl font-bold text-primary">
                      {formatCurrency(calculationResult.total_cost)}
                    </p>
                  </div>
                  <div className="space-y-2">
                    <p className="text-sm text-muted-foreground">單位成本</p>
                    <p className="text-2xl font-bold">
                      {formatCurrency(calculationResult.unit_cost)}
                    </p>
                  </div>
                  <div className="space-y-2">
                    <p className="text-sm text-muted-foreground">建議售價</p>
                    <p className="text-2xl font-bold text-green-600">
                      {formatCurrency(calculationResult.final_price)}
                    </p>
                  </div>
                </div>
                
                {calculationResult.details && calculationResult.details.length > 0 && (
                  <div className="mt-6">
                    <h4 className="font-semibold mb-2">成本明細</h4>
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead>項目</TableHead>
                          <TableHead>說明</TableHead>
                          <TableHead className="text-right">單價</TableHead>
                          <TableHead className="text-right">數量</TableHead>
                          <TableHead className="text-right">總價</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {calculationResult.details.map((detail: any, index: number) => (
                          <TableRow key={index}>
                            <TableCell>{detail.type}</TableCell>
                            <TableCell>{detail.description}</TableCell>
                            <TableCell className="text-right">
                              {formatCurrency(detail.unit_cost)}
                            </TableCell>
                            <TableCell className="text-right">{detail.quantity}</TableCell>
                            <TableCell className="text-right">
                              {formatCurrency(detail.total_cost)}
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </div>
                )}
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="templates" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>成本模板</CardTitle>
              <Button>新增模板</Button>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>模板名稱</TableHead>
                    <TableHead>描述</TableHead>
                    <TableHead>製程類型</TableHead>
                    <TableHead>類別</TableHead>
                    <TableHead className="text-right">總成本</TableHead>
                    <TableHead>狀態</TableHead>
                    <TableHead>操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {templates.map((template) => (
                    <TableRow key={template.id}>
                      <TableCell className="font-medium">{template.name}</TableCell>
                      <TableCell>{template.description}</TableCell>
                      <TableCell>{template.process_type}</TableCell>
                      <TableCell>{template.category}</TableCell>
                      <TableCell className="text-right">
                        {formatCurrency(template.total_cost)}
                      </TableCell>
                      <TableCell>
                        <span className={`px-2 py-1 rounded text-xs ${
                          template.is_active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                        }`}>
                          {template.is_active ? '啟用' : '停用'}
                        </span>
                      </TableCell>
                      <TableCell>
                        <Button size="sm" variant="outline">編輯</Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="materials" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>材料成本</CardTitle>
              <Button>新增材料</Button>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>材料名稱</TableHead>
                    <TableHead>類型</TableHead>
                    <TableHead>規格</TableHead>
                    <TableHead className="text-right">單價</TableHead>
                    <TableHead>單位</TableHead>
                    <TableHead>供應商</TableHead>
                    <TableHead>操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {materials.map((material) => (
                    <TableRow key={material.id}>
                      <TableCell className="font-medium">{material.name}</TableCell>
                      <TableCell>{material.type}</TableCell>
                      <TableCell>{material.specification}</TableCell>
                      <TableCell className="text-right">
                        {formatCurrency(material.unit_price, material.currency)}
                      </TableCell>
                      <TableCell>{material.unit}</TableCell>
                      <TableCell>{material.supplier}</TableCell>
                      <TableCell>
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => {
                            setSelectedMaterial(material);
                            setShowMaterialDialog(true);
                          }}
                        >
                          編輯
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="rates" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>加工費率</CardTitle>
              <Button>新增費率</Button>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>製程類型</TableHead>
                    <TableHead>設備名稱</TableHead>
                    <TableHead className="text-right">時薪費率</TableHead>
                    <TableHead className="text-right">設置成本</TableHead>
                    <TableHead className="text-right">最低收費</TableHead>
                    <TableHead>幣別</TableHead>
                    <TableHead>操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {processingRates.map((rate) => (
                    <TableRow key={rate.id}>
                      <TableCell className="font-medium">{rate.process_type}</TableCell>
                      <TableCell>{rate.equipment_name || '-'}</TableCell>
                      <TableCell className="text-right">
                        {formatCurrency(rate.hourly_rate, rate.currency)}
                      </TableCell>
                      <TableCell className="text-right">
                        {formatCurrency(rate.setup_cost, rate.currency)}
                      </TableCell>
                      <TableCell className="text-right">
                        {formatCurrency(rate.minimum_charge, rate.currency)}
                      </TableCell>
                      <TableCell>{rate.currency}</TableCell>
                      <TableCell>
                        <Button size="sm" variant="outline">編輯</Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {/* 材料編輯對話框 */}
      <Dialog open={showMaterialDialog} onOpenChange={setShowMaterialDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>編輯材料成本</DialogTitle>
          </DialogHeader>
          {selectedMaterial && (
            <div className="space-y-4">
              <div>
                <Label>材料名稱</Label>
                <Input
                  value={selectedMaterial.name}
                  onChange={(e) => setSelectedMaterial({
                    ...selectedMaterial,
                    name: e.target.value
                  })}
                />
              </div>
              <div>
                <Label>單價</Label>
                <Input
                  type="number"
                  value={selectedMaterial.unit_price}
                  onChange={(e) => setSelectedMaterial({
                    ...selectedMaterial,
                    unit_price: parseFloat(e.target.value) || 0
                  })}
                />
              </div>
              <div>
                <Label>變更原因</Label>
                <Input
                  placeholder="請輸入價格變更原因"
                  onChange={(e) => setSelectedMaterial({
                    ...selectedMaterial,
                    price_change_reason: e.target.value
                  })}
                />
              </div>
            </div>
          )}
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowMaterialDialog(false)}>
              取消
            </Button>
            <Button onClick={handleUpdateMaterial}>
              確認更新
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}