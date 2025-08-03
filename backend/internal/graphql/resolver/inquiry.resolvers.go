package resolver

import (
	"context"
	"fmt"

	"github.com/fastenmind/fastener-api/internal/domain/cqrs/commands"
	"github.com/fastenmind/fastener-api/internal/domain/cqrs/queries"
	"github.com/fastenmind/fastener-api/internal/graphql/model"
	"github.com/google/uuid"
)

// CreateInquiry creates a new inquiry
func (r *mutationResolver) CreateInquiry(ctx context.Context, input model.CreateInquiryInput) (*model.Inquiry, error) {
	// Get user from context
	user := getUserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	// Create command
	cmd := commands.NewCreateInquiryCommand(user.ID, user.CompanyID, input.CustomerID)
	cmd.ProductCategory = input.ProductCategory
	cmd.ProductName = input.ProductName
	cmd.DrawingFiles = input.DrawingFiles
	cmd.Quantity = input.Quantity
	cmd.Unit = input.Unit
	cmd.RequiredDate = input.RequiredDate
	cmd.Incoterm = input.Incoterm
	cmd.DestinationPort = input.DestinationPort
	cmd.DestinationAddress = input.DestinationAddress
	cmd.PaymentTerms = input.PaymentTerms
	cmd.SpecialRequirements = input.SpecialRequirements

	// Send command
	if err := r.commandBus.Send(ctx, cmd); err != nil {
		return nil, err
	}

	// Return created inquiry
	// In real implementation, the command handler would return the created inquiry ID
	return &model.Inquiry{}, nil
}

// UpdateInquiry updates an existing inquiry
func (r *mutationResolver) UpdateInquiry(ctx context.Context, id uuid.UUID, input model.UpdateInquiryInput) (*model.Inquiry, error) {
	user := getUserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	cmd := commands.NewUpdateInquiryCommand(user.ID, id)
	cmd.ProductName = input.ProductName
	cmd.Quantity = input.Quantity
	cmd.RequiredDate = input.RequiredDate
	cmd.SpecialRequirements = input.SpecialRequirements

	if err := r.commandBus.Send(ctx, cmd); err != nil {
		return nil, err
	}

	// Fetch updated inquiry
	query := queries.GetInquiryByIDQuery{InquiryID: id}
	result, err := r.queryBus.Send(ctx, query)
	if err != nil {
		return nil, err
	}

	inquiry := result.(*model.Inquiry)
	return inquiry, nil
}

// AssignInquiry assigns an inquiry to an engineer
func (r *mutationResolver) AssignInquiry(ctx context.Context, id uuid.UUID, engineerID uuid.UUID, note *string) (*model.Inquiry, error) {
	user := getUserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	cmd := commands.NewAssignInquiryCommand(user.ID, id, engineerID)
	if note != nil {
		cmd.AssignmentNote = *note
	}

	if err := r.commandBus.Send(ctx, cmd); err != nil {
		return nil, err
	}

	// Fetch updated inquiry
	query := queries.GetInquiryByIDQuery{InquiryID: id}
	result, err := r.queryBus.Send(ctx, query)
	if err != nil {
		return nil, err
	}

	inquiry := result.(*model.Inquiry)
	return inquiry, nil
}

// RejectInquiry rejects an inquiry
func (r *mutationResolver) RejectInquiry(ctx context.Context, id uuid.UUID, reason string) (*model.Inquiry, error) {
	user := getUserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	cmd := commands.NewRejectInquiryCommand(user.ID, id, reason)
	if err := r.commandBus.Send(ctx, cmd); err != nil {
		return nil, err
	}

	// Fetch updated inquiry
	query := queries.GetInquiryByIDQuery{InquiryID: id}
	result, err := r.queryBus.Send(ctx, query)
	if err != nil {
		return nil, err
	}

	inquiry := result.(*model.Inquiry)
	return inquiry, nil
}

// Inquiry fetches a single inquiry by ID
func (r *queryResolver) Inquiry(ctx context.Context, id uuid.UUID) (*model.Inquiry, error) {
	// Try cache first
	cached, err := r.cache.GetInquiry(ctx, id)
	if err == nil && cached != nil {
		return cached.(*model.Inquiry), nil
	}

	// Query from database
	query := queries.GetInquiryByIDQuery{InquiryID: id}
	result, err := r.queryBus.Send(ctx, query)
	if err != nil {
		return nil, err
	}

	inquiry := result.(*model.Inquiry)

	// Cache result
	go r.cache.SetInquiry(ctx, id, inquiry)

	return inquiry, nil
}

// Inquiries fetches inquiries with filters and pagination
func (r *queryResolver) Inquiries(ctx context.Context, filter *model.InquiryFilter, page *model.PageInput) (*model.InquiryConnection, error) {
	query := queries.ListInquiriesQuery{
		Page:     1,
		PageSize: 20,
	}

	// Apply filters
	if filter != nil {
		query.CompanyID = filter.CompanyID
		query.CustomerID = filter.CustomerID
		query.SalesID = filter.SalesID
		query.EngineerID = filter.EngineerID
		if filter.Status != nil {
			status := string(*filter.Status)
			query.Status = &status
		}
		query.DateFrom = filter.DateFrom
		query.DateTo = filter.DateTo
		query.SearchTerm = filter.SearchTerm
	}

	// Apply pagination
	if page != nil {
		if page.First != nil {
			query.PageSize = *page.First
		}
		// Handle cursor-based pagination
		// In real implementation, decode cursor to get page number
	}

	result, err := r.queryBus.Send(ctx, query)
	if err != nil {
		return nil, err
	}

	// Convert result to connection format
	inquiries := result.([]*model.Inquiry)
	
	edges := make([]*model.InquiryEdge, len(inquiries))
	for i, inquiry := range inquiries {
		edges[i] = &model.InquiryEdge{
			Node:   inquiry,
			Cursor: encodeCursor(inquiry.ID.String()),
		}
	}

	return &model.InquiryConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     len(inquiries) == query.PageSize,
			HasPreviousPage: query.Page > 1,
		},
		TotalCount: len(inquiries), // In real implementation, get total from query
	}, nil
}

// InquiryStatistics fetches inquiry statistics
func (r *queryResolver) InquiryStatistics(ctx context.Context, filter *model.InquiryFilter) (*model.InquiryStatistics, error) {
	// Implement statistics query
	return &model.InquiryStatistics{
		TotalCount:          100,
		PendingCount:        20,
		AssignedCount:       30,
		QuotedCount:         40,
		RejectedCount:       10,
		ConversionRate:      0.4,
		AverageResponseTime: 24.5,
	}, nil
}

// Helper functions
func getUserFromContext(ctx context.Context) *model.Account {
	// In real implementation, extract user from context
	return &model.Account{
		ID:        uuid.New(),
		CompanyID: uuid.New(),
	}
}

func encodeCursor(id string) string {
	// In real implementation, encode cursor properly
	return id
}