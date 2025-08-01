'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import assignmentService from '@/services/assignment.service';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent } from '@/components/ui/card';
import { useToast } from '@/components/ui/use-toast';
import { UserPlus, Loader2, Zap, Users, AlertCircle } from 'lucide-react';

interface InquiryAssignmentProps {
  inquiryId: string;
  inquiryNo: string;
  productCategory: string;
  currentEngineer?: {
    id: string;
    full_name: string;
  };
  onAssignmentComplete?: () => void;
}

export default function InquiryAssignment({
  inquiryId,
  inquiryNo,
  productCategory,
  currentEngineer,
  onAssignmentComplete,
}: InquiryAssignmentProps) {
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedEngineerId, setSelectedEngineerId] = useState('');
  const [assignmentReason, setAssignmentReason] = useState('');
  const [assignmentType, setAssignmentType] = useState<'manual' | 'reassign'>('manual');

  // 獲取可用工程師
  const { data: availableEngineers = [], isLoading: loadingEngineers } = useQuery({
    queryKey: ['available-engineers', productCategory],
    queryFn: () => assignmentService.getAvailableEngineers(productCategory),
    enabled: isDialogOpen,
  });

  // 獲取建議工程師
  const { data: suggestion, isLoading: loadingSuggestion } = useQuery({
    queryKey: ['engineer-suggestion', inquiryId],
    queryFn: () => assignmentService.suggestEngineer(inquiryId),
    enabled: isDialogOpen,
  });

  // 自動分派
  const autoAssignMutation = useMutation({
    mutationFn: () => assignmentService.autoAssign(inquiryId),
    onSuccess: () => {
      toast({
        title: '自動分派成功',
        description: '詢價單已根據規則自動分派給工程師',
      });
      queryClient.invalidateQueries({ queryKey: ['inquiries'] });
      onAssignmentComplete?.();
    },
    onError: (error: any) => {
      toast({
        title: '自動分派失敗',
        description: error.message || '無法完成自動分派',
        variant: 'destructive',
      });
    },
  });

  // 手動分派
  const manualAssignMutation = useMutation({
    mutationFn: (data: { engineer_id: string; reason?: string; type: 'manual' | 'reassign' }) =>
      assignmentService.manualAssign({
        inquiry_id: inquiryId,
        engineer_id: data.engineer_id,
        reason: data.reason,
        assignment_type: data.type,
      }),
    onSuccess: () => {
      toast({
        title: '分派成功',
        description: '詢價單已成功分派給指定工程師',
      });
      setIsDialogOpen(false);
      queryClient.invalidateQueries({ queryKey: ['inquiries'] });
      onAssignmentComplete?.();
    },
    onError: (error: any) => {
      toast({
        title: '分派失敗',
        description: error.message || '無法完成分派',
        variant: 'destructive',
      });
    },
  });

  const handleManualAssign = () => {
    if (!selectedEngineerId) {
      toast({
        title: '請選擇工程師',
        variant: 'destructive',
      });
      return;
    }

    manualAssignMutation.mutate({
      engineer_id: selectedEngineerId,
      reason: assignmentReason,
      type: assignmentType,
    });
  };

  const handleOpenDialog = () => {
    setAssignmentType(currentEngineer ? 'reassign' : 'manual');
    setSelectedEngineerId('');
    setAssignmentReason('');
    setIsDialogOpen(true);
  };

  return (
    <>
      <div className="flex gap-2">
        {!currentEngineer && (
          <Button
            variant="outline"
            size="sm"
            onClick={() => autoAssignMutation.mutate()}
            disabled={autoAssignMutation.isPending}
          >
            {autoAssignMutation.isPending ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <Zap className="mr-2 h-4 w-4" />
            )}
            自動分派
          </Button>
        )}
        
        <Button
          variant={currentEngineer ? 'outline' : 'default'}
          size="sm"
          onClick={handleOpenDialog}
        >
          <UserPlus className="mr-2 h-4 w-4" />
          {currentEngineer ? '重新分派' : '手動分派'}
        </Button>
      </div>

      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>
              {currentEngineer ? '重新分派詢價單' : '手動分派詢價單'}
            </DialogTitle>
          </DialogHeader>

          <div className="space-y-4">
            {/* 詢價單資訊 */}
            <Card>
              <CardContent className="pt-4">
                <div className="grid grid-cols-2 gap-2 text-sm">
                  <div>
                    <span className="text-gray-500">詢價單號：</span>
                    <span className="font-medium">{inquiryNo}</span>
                  </div>
                  <div>
                    <span className="text-gray-500">產品類別：</span>
                    <Badge variant="outline">{productCategory}</Badge>
                  </div>
                  {currentEngineer && (
                    <div className="col-span-2">
                      <span className="text-gray-500">目前負責：</span>
                      <span className="font-medium ml-1">{currentEngineer.full_name}</span>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* 建議工程師 */}
            {loadingSuggestion ? (
              <div className="text-center py-4">
                <Loader2 className="h-6 w-6 animate-spin mx-auto" />
              </div>
            ) : suggestion?.suggested_engineer ? (
              <Card className="border-blue-200 bg-blue-50">
                <CardContent className="pt-4">
                  <div className="flex items-start gap-2">
                    <AlertCircle className="h-5 w-5 text-blue-600 mt-0.5" />
                    <div className="flex-1">
                      <p className="font-medium text-blue-900">系統建議</p>
                      <p className="text-sm text-blue-700 mt-1">
                        建議分派給 <span className="font-semibold">{suggestion.suggested_engineer.engineer_name}</span>
                        {suggestion.reason && ` - ${suggestion.reason}`}
                      </p>
                      {suggestion.matching_rules.length > 0 && (
                        <div className="mt-2">
                          <p className="text-xs text-blue-600">匹配規則：</p>
                          {suggestion.matching_rules.map((rule) => (
                            <Badge key={rule.id} variant="secondary" className="text-xs mr-1">
                              {rule.rule_name}
                            </Badge>
                          ))}
                        </div>
                      )}
                    </div>
                  </div>
                </CardContent>
              </Card>
            ) : null}

            {/* 工程師選擇 */}
            <div className="space-y-2">
              <Label>選擇工程師</Label>
              <Select value={selectedEngineerId} onValueChange={setSelectedEngineerId}>
                <SelectTrigger>
                  <SelectValue placeholder="選擇要分派的工程師" />
                </SelectTrigger>
                <SelectContent>
                  {loadingEngineers ? (
                    <div className="py-2 text-center text-sm text-gray-500">載入中...</div>
                  ) : availableEngineers.length > 0 ? (
                    availableEngineers.map((engineer) => (
                      <SelectItem key={engineer.engineer_id} value={engineer.engineer_id}>
                        <div className="flex items-center justify-between w-full">
                          <span>{engineer.engineer_name}</span>
                          <div className="flex items-center gap-2 text-xs text-gray-500">
                            <span>進行中: {engineer.current_inquiries}</span>
                            {suggestion?.suggested_engineer?.engineer_id === engineer.engineer_id && (
                              <Badge variant="secondary" className="text-xs">建議</Badge>
                            )}
                          </div>
                        </div>
                      </SelectItem>
                    ))
                  ) : (
                    <div className="py-2 text-center text-sm text-gray-500">
                      沒有可處理 {productCategory} 類產品的工程師
                    </div>
                  )}
                </SelectContent>
              </Select>
            </div>

            {/* 分派原因 */}
            <div className="space-y-2">
              <Label htmlFor="reason">分派原因（選填）</Label>
              <Textarea
                id="reason"
                placeholder="請輸入分派或重新分派的原因..."
                value={assignmentReason}
                onChange={(e) => setAssignmentReason(e.target.value)}
                rows={3}
              />
            </div>

            {/* 工程師工作量預覽 */}
            {selectedEngineerId && (
              <Card>
                <CardContent className="pt-4">
                  <p className="text-sm font-medium mb-2">工程師工作量</p>
                  {(() => {
                    const engineer = availableEngineers.find(e => e.engineer_id === selectedEngineerId);
                    if (!engineer) return null;
                    
                    return (
                      <div className="grid grid-cols-2 gap-2 text-sm">
                        <div>
                          <span className="text-gray-500">目前進行中：</span>
                          <span className="font-medium ml-1">{engineer.current_inquiries} 件</span>
                        </div>
                        <div>
                          <span className="text-gray-500">今日完成：</span>
                          <span className="font-medium ml-1">{engineer.completed_today} 件</span>
                        </div>
                        <div className="col-span-2">
                          <span className="text-gray-500">專長領域：</span>
                          {engineer.skill_categories.map((cat) => (
                            <Badge key={cat} variant="outline" className="text-xs ml-1">
                              {cat}
                            </Badge>
                          ))}
                        </div>
                      </div>
                    );
                  })()}
                </CardContent>
              </Card>
            )}
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsDialogOpen(false)}>
              取消
            </Button>
            <Button 
              onClick={handleManualAssign}
              disabled={!selectedEngineerId || manualAssignMutation.isPending}
            >
              {manualAssignMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              確認分派
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}