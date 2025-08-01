export interface Quote {
  id: string;
  quote_no: string;
  inquiry_id?: string;
  customer_id: string;
  status: 'draft' | 'pending_approval' | 'approved' | 'rejected' | 'sent' | 'expired';
  approval_status?: string;
  validity_days: number;
  valid_until?: string;
  payment_terms?: string;
  delivery_terms?: string;
  total_amount: number;
  approved_amount?: number;
  remarks?: string;
  created_by?: string;
  updated_by?: string;
  approved_by?: string;
  approved_at?: string;
  current_version_id?: string;
  template_id?: string;
  created_at: string;
  updated_at: string;
  
  // Relations
  customer?: Customer;
  inquiry?: Inquiry;
  created_by_user?: User;
  updated_by_user?: User;
  approved_by_user?: User;
  template?: QuoteTemplate;
}

export interface QuoteVersion {
  id: string;
  quote_id: string;
  version_number: number;
  version_notes?: string;
  is_current: boolean;
  created_by: string;
  created_at: string;
  
  // Relations
  items?: QuoteItem[];
  terms?: QuoteTerm[];
  creator?: User;
}

export interface QuoteItem {
  id: string;
  quote_version_id: string;
  item_no: number;
  product_name: string;
  product_specs?: string;
  quantity: number;
  unit: string;
  unit_price: number;
  total_price: number;
  cost_calculation_id?: string;
  margin_percentage?: number;
  notes?: string;
  created_at: string;
  
  // Relations
  cost_calculation?: CostCalculation;
}

export interface QuoteTerm {
  id: string;
  quote_version_id: string;
  term_type: string;
  term_content: string;
  sort_order: number;
  created_at: string;
}

export interface QuoteApproval {
  id: string;
  quote_id: string;
  quote_version_id: string;
  approval_level: number;
  approver_role: string;
  required_approver_id?: string;
  actual_approver_id?: string;
  approval_status: 'pending' | 'approved' | 'rejected';
  approval_notes?: string;
  approved_at?: string;
  created_at: string;
  
  // Relations
  required_approver?: User;
  actual_approver?: User;
}

export interface QuoteActivityLog {
  id: string;
  quote_id: string;
  quote_version_id?: string;
  activity_type: string;
  activity_description?: string;
  activity_data?: any;
  performed_by: string;
  performed_at: string;
  
  // Relations
  performer?: User;
}

export interface QuoteTermsTemplate {
  id: string;
  template_name: string;
  template_type: string;
  content: string;
  language: string;
  is_default: boolean;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface QuoteTemplate {
  id: string;
  template_name: string;
  description?: string;
  header_logo_path?: string;
  footer_content?: string;
  terms_conditions?: string;
  css_styles?: string;
  is_default: boolean;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

// Request types
export interface CreateQuoteRequest {
  inquiry_id: string;
  customer_id: string;
  validity_days?: number;
  payment_terms?: string;
  delivery_terms?: string;
  remarks?: string;
  items: QuoteItemRequest[];
  terms?: QuoteTermRequest[];
  use_template?: boolean;
  template_id?: string;
}

export interface QuoteItemRequest {
  product_name: string;
  product_specs?: string;
  quantity: number;
  unit: string;
  unit_price: number;
  cost_calculation_id?: string;
  notes?: string;
}

export interface QuoteTermRequest {
  term_type: string;
  term_content: string;
  sort_order?: number;
}

export interface UpdateQuoteRequest {
  create_new_version?: boolean;
  version_notes?: string;
  validity_days?: number;
  payment_terms?: string;
  delivery_terms?: string;
  remarks?: string;
  items?: QuoteItemRequest[];
  terms?: QuoteTermRequest[];
}

export interface SubmitApprovalRequest {
  notes?: string;
}

export interface ApproveQuoteRequest {
  approved: boolean;
  notes?: string;
}

export interface SendQuoteRequest {
  recipient_email: string;
  recipient_name?: string;
  cc_emails?: string[];
  subject?: string;
  message?: string;
  attach_pdf?: boolean;
  attachment_ids?: string[];
}

// Import related types
interface Customer {
  id: string;
  customer_code: string;
  name: string;
  name_en?: string;
  email?: string;
  phone?: string;
}

interface Inquiry {
  id: string;
  inquiry_no: string;
  customer_id: string;
  status: string;
}

interface User {
  id: string;
  username: string;
  name: string;
  email: string;
  role: string;
}

interface CostCalculation {
  id: string;
  calculation_name: string;
  total_cost: number;
  unit_price: number;
}