package gopb

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/shopspring/decimal"
	"github.com/xefino/protobuf-gen-go/utils"
	"gopkg.in/yaml.v3"
)

// ProviderAlternates contains alternate values for the Provider enum
var ProviderAlternates = map[string]Provider{
	"":        Provider_None,
	"polygon": Provider_Polygon,
}

// ProviderMapping contains alternate names for the Provider enum
var ProviderMapping = map[Provider]string{
	Provider_None:    "",
	Provider_Polygon: "polygon",
}

// AssetClassAlternates contains alternative values for the Financial.Common.AssetClass enum
var AssetClassAlternates = map[string]Financial_Common_AssetClass{
	"":                 utils.NoValue[Financial_Common_AssetClass](),
	"stocks":           Financial_Common_Stock,
	"options":          Financial_Common_Option,
	"crypto":           Financial_Common_Crypto,
	"fx":               Financial_Common_ForeignExchange,
	"Foreign Exchange": Financial_Common_ForeignExchange,
	"otc":              Financial_Common_OverTheCounter,
	"OTC":              Financial_Common_OverTheCounter,
	"indices":          Financial_Common_Indices,
	"index":            Financial_Common_Indices,
	"Index":            Financial_Common_Indices,
}

// AssetClassMapping contains alternate names for the Financial.Common.AssetClass enum
var AssetClassMapping = map[Financial_Common_AssetClass]string{
	Financial_Common_ForeignExchange: "Foreign Exchange",
	Financial_Common_OverTheCounter:  "OTC",
}

// AssetTypeAlternates contains alternative values for the Financial.Common.AssetType enum
var AssetTypeAlternates = map[string]Financial_Common_AssetType{
	"":                        utils.NoValue[Financial_Common_AssetType](),
	"CS":                      Financial_Common_CommonShare,
	"Common Share":            Financial_Common_CommonShare,
	"OS":                      Financial_Common_OrdinaryShare,
	"Ordinary Share":          Financial_Common_OrdinaryShare,
	"NYRS":                    Financial_Common_NewYorkRegistryShares,
	"New York Registry Share": Financial_Common_NewYorkRegistryShares,
	"ADRC":                    Financial_Common_AmericanDepositoryReceiptCommon,
	"Common ADR":              Financial_Common_AmericanDepositoryReceiptCommon,
	"ADRP":                    Financial_Common_AmericanDepositoryReceiptPreferred,
	"Preferred ADR":           Financial_Common_AmericanDepositoryReceiptPreferred,
	"ADRR":                    Financial_Common_AmericanDepositoryReceiptRights,
	"ADR Right":               Financial_Common_AmericanDepositoryReceiptRights,
	"ADRW":                    Financial_Common_AmericanDepositoryReceiptWarrants,
	"ADR Warrant":             Financial_Common_AmericanDepositoryReceiptWarrants,
	"GDR":                     Financial_Common_GlobalDepositoryReceipts,
	"UNIT":                    Financial_Common_Unit,
	"RIGHT":                   Financial_Common_Rights,
	"Right":                   Financial_Common_Rights,
	"PFD":                     Financial_Common_PreferredStock,
	"Preferred Stock":         Financial_Common_PreferredStock,
	"FUND":                    Financial_Common_Fund,
	"SP":                      Financial_Common_StructuredProduct,
	"Structured Product":      Financial_Common_StructuredProduct,
	"WARRANT":                 Financial_Common_Warrant,
	"INDEX":                   Financial_Common_Index,
	"ETF":                     Financial_Common_ExchangeTradedFund,
	"ETN":                     Financial_Common_ExchangeTradedNote,
	"ETV":                     Financial_Common_ExchangeTradeVehicle,
	"ETS":                     Financial_Common_SingleSecurityETF,
	"BOND":                    Financial_Common_CorporateBond,
	"Corporate Bond":          Financial_Common_CorporateBond,
	"AGEN":                    Financial_Common_AgencyBond,
	"Agency Bond":             Financial_Common_AgencyBond,
	"EQLK":                    Financial_Common_EquityLinkedBond,
	"Equity-Linked Bond":      Financial_Common_EquityLinkedBond,
	"BASKET":                  Financial_Common_Basket,
	"LT":                      Financial_Common_LiquidatingTrust,
	"Liquidating Trust":       Financial_Common_LiquidatingTrust,
	"OTHER":                   Financial_Common_Others,
	"Other":                   Financial_Common_Others,
	"None":                    Financial_Common_None,
}

// AssetTypeMapping contains alternate names for the Financial.Common.AssetType enum
var AssetTypeMapping = map[Financial_Common_AssetType]string{
	Financial_Common_CommonShare:                        "Common Share",
	Financial_Common_OrdinaryShare:                      "Ordinary Share",
	Financial_Common_NewYorkRegistryShares:              "New York Registry Share",
	Financial_Common_AmericanDepositoryReceiptCommon:    "Common ADR",
	Financial_Common_AmericanDepositoryReceiptPreferred: "Preferred ADR",
	Financial_Common_AmericanDepositoryReceiptRights:    "ADR Right",
	Financial_Common_AmericanDepositoryReceiptWarrants:  "ADR Warrant",
	Financial_Common_GlobalDepositoryReceipts:           "GDR",
	Financial_Common_Rights:                             "Right",
	Financial_Common_PreferredStock:                     "Preferred Stock",
	Financial_Common_StructuredProduct:                  "Structured Product",
	Financial_Common_ExchangeTradedFund:                 "ETF",
	Financial_Common_ExchangeTradedNote:                 "ETN",
	Financial_Common_ExchangeTradeVehicle:               "ETV",
	Financial_Common_SingleSecurityETF:                  "ETS",
	Financial_Common_CorporateBond:                      "Corporate Bond",
	Financial_Common_AgencyBond:                         "Agency Bond",
	Financial_Common_EquityLinkedBond:                   "Equity-Linked Bond",
	Financial_Common_LiquidatingTrust:                   "Liquidating Trust",
	Financial_Common_Others:                             "Other",
	Financial_Common_None:                               "",
}

// LocalAlternates contains alternative values for the Financial.Common.Locale enum
var LocaleAlternates = map[string]Financial_Common_Locale{
	"":       utils.NoValue[Financial_Common_Locale](),
	"us":     Financial_Common_US,
	"global": Financial_Common_Global,
}

// ExchangeTypeAlternates contains alternative values for the Financial.Dividends.Frequency enum
var DividendFrequencyAlternates = map[string]Financial_Dividends_Frequency{
	"None": Financial_Dividends_NoFrequency,
	"":     Financial_Dividends_NoFrequency,
}

// DividendFrequencyMapping contains alternate names for the Financial.Dividends.Frequency enum
var DividendFrequencyMapping = map[Financial_Dividends_Frequency]string{
	Financial_Dividends_NoFrequency: "",
}

// ExchangeTypeAlternates contains alternative values for the Financial.Exchanges.Type enum
var ExchangeTypeAlternates = map[string]Financial_Exchanges_Type{
	"exchange": Financial_Exchanges_Exchange,
}

// OptionContractTypeAlternates contains alternative values for the Financial.Options.ContractType enum
var OptionContractTypeAlternates = map[string]Financial_Options_ContractType{
	"call":  Financial_Options_Call,
	"put":   Financial_Options_Put,
	"other": Financial_Options_Other,
}

// OptionExerciseStyleAlternates contains alternative values for the Financial.Options.ExerciseStyle enum
var OptionExerciseStyleAlternates = map[string]Financial_Options_ExerciseStyle{
	"american": Financial_Options_American,
	"european": Financial_Options_European,
	"bermudan": Financial_Options_Bermudan,
}

// UnderlyingTypeAlternates contains alternative values for the Financial.Options.UnderlyingType enum
var OptionUnderlyingTypeAlternates = map[string]Financial_Options_UnderlyingType{
	"equity":   Financial_Options_Equity,
	"currency": Financial_Options_Currency,
}

// QuoteConditionAlternates contains alternative values for the Financial.Quotes.Condition enum
var QuoteConditionAlternates = map[string]Financial_Quotes_Condition{
	"-1":                                Financial_Quotes_Invalid,
	"Regular, Two-Sided Open":           Financial_Quotes_RegularTwoSidedOpen,
	"Regular, One-Sided Open":           Financial_Quotes_RegularOneSidedOpen,
	"Slow Ask":                          Financial_Quotes_SlowAsk,
	"Slow Bid":                          Financial_Quotes_SlowBid,
	"Slow Bid, Ask":                     Financial_Quotes_SlowBidAsk,
	"Slow Due, LRP Bid":                 Financial_Quotes_SlowDueLRPBid,
	"Slow Due, LRP Ask":                 Financial_Quotes_SlowDueLRPAsk,
	"Slow Due, NYSE LRP":                Financial_Quotes_SlowDueNYSELRP,
	"Slow Due Set, Slow List, Bid, Ask": Financial_Quotes_SlowDueSetSlowListBidAsk,
	"Manual Ask, Automated Bid":         Financial_Quotes_ManualAskAutomatedBid,
	"Manual Bid, Automated Ask":         Financial_Quotes_ManualBidAutomatedAsk,
	"Manual Bid and Ask":                Financial_Quotes_ManualBidAndAsk,
	"Fast Trading":                      Financial_Quotes_FastTrading,
	"Tading Range Indicated":            Financial_Quotes_TradingRangeIndicated,
	"Market-Maker Quotes Closed":        Financial_Quotes_MarketMakerQuotesClosed,
	"Non-Firm":                          Financial_Quotes_NonFirm,
	"News Dissemination":                Financial_Quotes_NewsDissemination,
	"Order Influx":                      Financial_Quotes_OrderInflux,
	"Order Imbalance":                   Financial_Quotes_OrderImbalance,
	"Due to Related Security, News Dissemination":    Financial_Quotes_DueToRelatedSecurityNewsDissemination,
	"Due to Related Security, News Pending":          Financial_Quotes_DueToRelatedSecurityNewsPending,
	"Additional Information":                         Financial_Quotes_AdditionalInformation,
	"News Pending":                                   Financial_Quotes_NewsPending,
	"Additional Information Due to Related Security": Financial_Quotes_AdditionalInformationDueToRelatedSecurity,
	"Due to Related Security":                        Financial_Quotes_DueToRelatedSecurity,
	"In View of Common":                              Financial_Quotes_InViewOfCommon,
	"Equipment Changeover":                           Financial_Quotes_EquipmentChangeover,
	"No Open, No Response":                           Financial_Quotes_NoOpenNoResponse,
	"Sub-Penny Trading":                              Financial_Quotes_SubPennyTrading,
	"Automated Bid; No Offer, No Bid":                Financial_Quotes_AutomatedBidNoOfferNoBid,
	"LULD Price Band":                                Financial_Quotes_LULDPriceBand,
	"Market-Wide Circuit Breaker, Level 1":           Financial_Quotes_MarketWideCircuitBreakerLevel1,
	"Market-Wide Circuit Breaker, Level 2":           Financial_Quotes_MarketWideCircuitBreakerLevel2,
	"Market-Wide Circuit Breaker, Level 3":           Financial_Quotes_MarketWideCircuitBreakerLevel3,
	"Republished LULD Price Band":                    Financial_Quotes_RepublishedLULDPriceBand,
	"On-Demand Auction":                              Financial_Quotes_OnDemandAuction,
	"Cash-Only Settlement":                           Financial_Quotes_CashOnlySettlement,
	"Next-Day Settlement":                            Financial_Quotes_NextDaySettlement,
	"LULD Trading Pause":                             Financial_Quotes_LULDTradingPause,
	"Slow Due LRP, Bid, Ask":                         Financial_Quotes_SlowDueLRPBidAsk,
	"Cancel":                                         Financial_Quotes_Cancel,
	"Corrected Price":                                Financial_Quotes_CorrectedPrice,
	"SIP-Generated":                                  Financial_Quotes_SIPGenerated,
	"Unknown":                                        Financial_Quotes_Unknown,
	"Crossed Market":                                 Financial_Quotes_CrossedMarket,
	"Locked Market":                                  Financial_Quotes_LockedMarket,
	"Depth on Offer Side":                            Financial_Quotes_DepthOnOfferSide,
	"Depth on Bid Side":                              Financial_Quotes_DepthOnBidSide,
	"Depth on Bid and Offer":                         Financial_Quotes_DepthOnBidAndOffer,
	"Pre-Opening Indication":                         Financial_Quotes_PreOpeningIndication,
	"Syndicate Bid":                                  Financial_Quotes_SyndicateBid,
	"Pre-Syndicate Bid":                              Financial_Quotes_PreSyndicateBid,
	"Penalty Bid":                                    Financial_Quotes_PenaltyBid,
	"CQS-Generated":                                  Financial_Quotes_CQSGenerated,
}

// QuoteConditionMapping contains alternate names for the Financial.Quotes.Condition enum
var QuoteConditionMapping = map[Financial_Quotes_Condition]string{
	Financial_Quotes_RegularTwoSidedOpen:                       "Regular, Two-Sided Open",
	Financial_Quotes_RegularOneSidedOpen:                       "Regular, One-Sided Open",
	Financial_Quotes_SlowAsk:                                   "Slow Ask",
	Financial_Quotes_SlowBid:                                   "Slow Bid",
	Financial_Quotes_SlowBidAsk:                                "Slow Bid, Ask",
	Financial_Quotes_SlowDueLRPBid:                             "Slow Due, LRP Bid",
	Financial_Quotes_SlowDueLRPAsk:                             "Slow Due, LRP Ask",
	Financial_Quotes_SlowDueNYSELRP:                            "Slow Due, NYSE LRP",
	Financial_Quotes_SlowDueSetSlowListBidAsk:                  "Slow Due Set, Slow List, Bid, Ask",
	Financial_Quotes_ManualAskAutomatedBid:                     "Manual Ask, Automated Bid",
	Financial_Quotes_ManualBidAutomatedAsk:                     "Manual Bid, Automated Ask",
	Financial_Quotes_ManualBidAndAsk:                           "Manual Bid and Ask",
	Financial_Quotes_FastTrading:                               "Fast Trading",
	Financial_Quotes_TradingRangeIndicated:                     "Tading Range Indicated",
	Financial_Quotes_MarketMakerQuotesClosed:                   "Market-Maker Quotes Closed",
	Financial_Quotes_NonFirm:                                   "Non-Firm",
	Financial_Quotes_NewsDissemination:                         "News Dissemination",
	Financial_Quotes_OrderInflux:                               "Order Influx",
	Financial_Quotes_OrderImbalance:                            "Order Imbalance",
	Financial_Quotes_DueToRelatedSecurityNewsDissemination:     "Due to Related Security, News Dissemination",
	Financial_Quotes_DueToRelatedSecurityNewsPending:           "Due to Related Security, News Pending",
	Financial_Quotes_AdditionalInformation:                     "Additional Information",
	Financial_Quotes_NewsPending:                               "News Pending",
	Financial_Quotes_AdditionalInformationDueToRelatedSecurity: "Additional Information Due to Related Security",
	Financial_Quotes_DueToRelatedSecurity:                      "Due to Related Security",
	Financial_Quotes_InViewOfCommon:                            "In View of Common",
	Financial_Quotes_EquipmentChangeover:                       "Equipment Changeover",
	Financial_Quotes_NoOpenNoResponse:                          "No Open, No Response",
	Financial_Quotes_SubPennyTrading:                           "Sub-Penny Trading",
	Financial_Quotes_AutomatedBidNoOfferNoBid:                  "Automated Bid; No Offer, No Bid",
	Financial_Quotes_LULDPriceBand:                             "LULD Price Band",
	Financial_Quotes_MarketWideCircuitBreakerLevel1:            "Market-Wide Circuit Breaker, Level 1",
	Financial_Quotes_MarketWideCircuitBreakerLevel2:            "Market-Wide Circuit Breaker, Level 2",
	Financial_Quotes_MarketWideCircuitBreakerLevel3:            "Market-Wide Circuit Breaker, Level 3",
	Financial_Quotes_RepublishedLULDPriceBand:                  "Republished LULD Price Band",
	Financial_Quotes_OnDemandAuction:                           "On-Demand Auction",
	Financial_Quotes_CashOnlySettlement:                        "Cash-Only Settlement",
	Financial_Quotes_NextDaySettlement:                         "Next-Day Settlement",
	Financial_Quotes_LULDTradingPause:                          "LULD Trading Pause",
	Financial_Quotes_SlowDueLRPBidAsk:                          "Slow Due LRP, Bid, Ask",
	Financial_Quotes_CorrectedPrice:                            "Corrected Price",
	Financial_Quotes_SIPGenerated:                              "SIP-Generated",
	Financial_Quotes_CrossedMarket:                             "Crossed Market",
	Financial_Quotes_LockedMarket:                              "Locked Market",
	Financial_Quotes_DepthOnOfferSide:                          "Depth on Offer Side",
	Financial_Quotes_DepthOnBidSide:                            "Depth on Bid Side",
	Financial_Quotes_DepthOnBidAndOffer:                        "Depth on Bid and Offer",
	Financial_Quotes_PreOpeningIndication:                      "Pre-Opening Indication",
	Financial_Quotes_SyndicateBid:                              "Syndicate Bid",
	Financial_Quotes_PreSyndicateBid:                           "Pre-Syndicate Bid",
	Financial_Quotes_PenaltyBid:                                "Penalty Bid",
	Financial_Quotes_CQSGenerated:                              "CQS-Generated",
}

// QuoteIndicatorAlternates contains alternative values for the Financial.Quotes.Indicator enum
var QuoteIndicatorAlternates = map[string]Financial_Quotes_Indicator{
	"NBB and/or NBO are Executable":                                 Financial_Quotes_NBBNBOExecutable,
	"NBB below Lower Band":                                          Financial_Quotes_NBBBelowLowerBand,
	"NBO above Upper Band":                                          Financial_Quotes_NBOAboveUpperBand,
	"NBB below Lower Band and NBO above Upper Band":                 Financial_Quotes_NBBBelowLowerBandAndNBOAboveUpperBand,
	"NBB equals Upper Band":                                         Financial_Quotes_NBBEqualsUpperBand,
	"NBO equals Lower Band":                                         Financial_Quotes_NBOEqualsLowerBand,
	"NBB equals Upper Band and NBO above Upper Band":                Financial_Quotes_NBBEqualsUpperBandAndNBOAboveUpperBand,
	"NBB below Lower Band and NBO equals Lower Band":                Financial_Quotes_NBBBelowLowerBandAndNBOEqualsLowerBand,
	"Bid Price above Upper Limit Price Band":                        Financial_Quotes_BidPriceAboveUpperLimitPriceBand,
	"Offer Price below Lower Limit Price Band":                      Financial_Quotes_OfferPriceBelowLowerLimitPriceBand,
	"Bid and Offer outside Price Band":                              Financial_Quotes_BidAndOfferOutsidePriceBand,
	"Opening Update":                                                Financial_Quotes_OpeningUpdate,
	"Intra-Day Update":                                              Financial_Quotes_IntraDayUpdate,
	"Restated Value":                                                Financial_Quotes_RestatedValue,
	"Suspended during Trading Halt or Trading Pause":                Financial_Quotes_SuspendedDuringTradingHalt,
	"Re-Opening Update":                                             Financial_Quotes_ReOpeningUpdate,
	"Outside Price Band Rule Hours":                                 Financial_Quotes_OutsidePriceBandRuleHours,
	"Auction Extension (Auction Collar Message)":                    Financial_Quotes_AuctionExtension,
	"LULD Price Band":                                               Financial_Quotes_LULDPriceBandInd,
	"Republished LULD Price Band":                                   Financial_Quotes_RepublishedLULDPriceBandInd,
	"NBB Limit State Entered":                                       Financial_Quotes_NBBLimitStateEntered,
	"NBB Limit State Exited":                                        Financial_Quotes_NBBLimitStateExited,
	"NBO Limit State Entered":                                       Financial_Quotes_NBOLimitStateEntered,
	"NBO Limit State Exited":                                        Financial_Quotes_NBOLimitStateExited,
	"NBB and NBO Limit State Entered":                               Financial_Quotes_NBBAndNBOLimitStateEntered,
	"NBB and NBO Limit State Exited":                                Financial_Quotes_NBBAndNBOLimitStateExited,
	"NBB Limit State Entered and NBO Limit State Exited":            Financial_Quotes_NBBLimitStateEnteredNBOLimitStateExited,
	"NBB Limit State Exited and NBO Limit State Entered":            Financial_Quotes_NBBLimitStateExitedNBOLimitStateEntered,
	"Deficient - Below Listing Requirements":                        Financial_Quotes_Deficient,
	"Delinquent - Late Filing":                                      Financial_Quotes_Delinquent,
	"Bankrupt and Deficient":                                        Financial_Quotes_BankruptAndDeficient,
	"Bankrupt and Delinquent":                                       Financial_Quotes_BankruptAndDelinquent,
	"Deficient and Delinquent":                                      Financial_Quotes_DeficientAndDelinquent,
	"Deficient, Delinquent, and Bankrupt":                           Financial_Quotes_DeficientDeliquentBankrupt,
	"Creations Suspended":                                           Financial_Quotes_CreationsSuspended,
	"Redemptions Suspended":                                         Financial_Quotes_RedemptionsSuspended,
	"Creations and/or Redemptions Suspended":                        Financial_Quotes_CreationsRedemptionsSuspended,
	"Normal Trading":                                                Financial_Quotes_NormalTrading,
	"Opening Delay":                                                 Financial_Quotes_OpeningDelay,
	"Trading Halt":                                                  Financial_Quotes_TradingHalt,
	"Resume":                                                        Financial_Quotes_TradingResume,
	"No Open / No Resume":                                           Financial_Quotes_NoOpenNoResume,
	"Price Indication":                                              Financial_Quotes_PriceIndication,
	"Trading Range Indication":                                      Financial_Quotes_TradingRangeIndication,
	"Market Imbalance Buy":                                          Financial_Quotes_MarketImbalanceBuy,
	"Market Imbalance Sell":                                         Financial_Quotes_MarketImbalanceSell,
	"Market On-Close Imbalance Buy":                                 Financial_Quotes_MarketOnCloseImbalanceBuy,
	"Market On Close Imbalance Sell":                                Financial_Quotes_MarketOnCloseImbalanceSell,
	"No Market Imbalance":                                           Financial_Quotes_NoMarketImbalance,
	"No Market, On-Close Imbalance":                                 Financial_Quotes_NoMarketOnCloseImbalance,
	"Short Sale Restriction":                                        Financial_Quotes_ShortSaleRestriction,
	"Limit Up-Limit Down":                                           Financial_Quotes_LimitUpLimitDown,
	"Quotation Resumption":                                          Financial_Quotes_QuotationResumption,
	"Trading Resumption":                                            Financial_Quotes_TradingResumption,
	"Volatility Trading Pause":                                      Financial_Quotes_VolatilityTradingPause,
	"Halt: News Pending":                                            Financial_Quotes_HaltNewsPending,
	"Update: News Dissemination":                                    Financial_Quotes_UpdateNewsDissemination,
	"Halt: Single Stock Trading Pause In Affect":                    Financial_Quotes_HaltSingleStockTradingPause,
	"Halt: Regulatory Extraordinary Market Activity":                Financial_Quotes_HaltRegulatoryExtraordinaryMarketActivity,
	"Halt: ETF":                                                     Financial_Quotes_HaltETF,
	"Halt: Information Requested":                                   Financial_Quotes_HaltInformationRequested,
	"Halt: Exchange Non-Compliance":                                 Financial_Quotes_HaltExchangeNonCompliance,
	"Halt: Filings Not Current":                                     Financial_Quotes_HaltFilingsNotCurrent,
	"Halt: SEC Trading Suspension":                                  Financial_Quotes_HaltSECTradingSuspension,
	"Halt: Regulatory Concern":                                      Financial_Quotes_HaltRegulatoryConcern,
	"Halt: Market Operations":                                       Financial_Quotes_HaltMarketOperations,
	"IPO Security: Not Yet Trading":                                 Financial_Quotes_IPOSecurityNotYetTrading,
	"Halt: Corporate Action":                                        Financial_Quotes_HaltCorporateAction,
	"Quotation Not Available":                                       Financial_Quotes_QuotationNotAvailable,
	"Halt: Volatility Trading Pause":                                Financial_Quotes_HaltVolatilityTradingPause,
	"Halt: Volatility Trading Pause - Straddle Condition":           Financial_Quotes_HaltVolatilityTradingPauseStraddleCondition,
	"Update: News and Resumption Times":                             Financial_Quotes_UpdateNewsAndResumptionTimes,
	"Halt: Single Stock Trading Pause - Quotes Only":                Financial_Quotes_HaltSingleStockTradingPauseQuotesOnly,
	"Resume: Qualification Issues Reviewed / Resolved":              Financial_Quotes_ResumeQualificationIssuesReviewedResolved,
	"Resume: Filing Requirements Satisfied / Resolved":              Financial_Quotes_ResumeFilingRequirementsSatisfiedResolved,
	"Resume: News Not Forthcoming":                                  Financial_Quotes_ResumeNewsNotForthcoming,
	"Resume: Qualifications - Maintenance Requirements Met":         Financial_Quotes_ResumeQualificationsMaintRequirementsMet,
	"Resume: Qualifications - Filings Met":                          Financial_Quotes_ResumeQualificationsFilingsMet,
	"Resume: Regulatory Auth":                                       Financial_Quotes_ResumeRegulatoryAuth,
	"New Issue Available":                                           Financial_Quotes_NewIssueAvailable,
	"Issue Available":                                               Financial_Quotes_IssueAvailable,
	"MWCB - Carry from Previous Day":                                Financial_Quotes_MWCBCarryFromPreviousDay,
	"MWCB - Resume":                                                 Financial_Quotes_MWCBResume,
	"IPO Security: Released for Quotation":                          Financial_Quotes_IPOSecurityReleasedForQuotation,
	"IPO Security: Positioning Window Extension":                    Financial_Quotes_IPOSecurityPositioningWindowExtension,
	"MWCB - Level 1":                                                Financial_Quotes_MWCBLevel1,
	"MWCB - Level 2":                                                Financial_Quotes_MWCBLevel2,
	"MWCB - Level 3":                                                Financial_Quotes_MWCBLevel3,
	"Halt: Sub-Penny Trading":                                       Financial_Quotes_HaltSubPennyTrading,
	"Order Imbalance":                                               Financial_Quotes_OrderImbalanceInd,
	"LULD Trading Paused":                                           Financial_Quotes_LULDTradingPaused,
	"Security Status: None":                                         Financial_Quotes_NONE,
	"Short Sales Restriction Activated":                             Financial_Quotes_ShortSalesRestrictionActivated,
	"Short Sales Restriction Continued":                             Financial_Quotes_ShortSalesRestrictionContinued,
	"Short Sales Restriction Deactivated":                           Financial_Quotes_ShortSalesRestrictionDeactivated,
	"Short Sales Restriction in Effect":                             Financial_Quotes_ShortSalesRestrictionInEffect,
	"Short Sales Restriction Max":                                   Financial_Quotes_ShortSalesRestrictionMax,
	"NBBO_NO_CHANGE":                                                Financial_Quotes_NBBONoChange,
	"NBBO: No Change":                                               Financial_Quotes_NBBONoChange,
	"NBBO_QUOTE_IS_NBBO":                                            Financial_Quotes_NBBOQuoteIsNBBO,
	"NBBO: Quote is NBBO":                                           Financial_Quotes_NBBOQuoteIsNBBO,
	"NBBO_NO_BB_NO_BO":                                              Financial_Quotes_NBBONoBBNoBO,
	"NBBO: No BB, No BO":                                            Financial_Quotes_NBBONoBBNoBO,
	"NBBO_BB_BO_SHORT_APPENDAGE":                                    Financial_Quotes_NBBOBBBOShortAppendage,
	"NBBO: BB / BO Short Appendage":                                 Financial_Quotes_NBBOBBBOShortAppendage,
	"NBBO_BB_BO_LONG_APPENDAGE":                                     Financial_Quotes_NBBOBBBOLongAppendage,
	"NBBO: BB / BO Long Appendage":                                  Financial_Quotes_NBBOBBBOLongAppendage,
	"HELD_TRADE_NOT_LAST_SALE_AND_NOT_ON_CONSOLIDATED":              Financial_Quotes_HeldTradeNotLastSaleNotConsolidated,
	"Held Trade not Last Sale, not Consolidated":                    Financial_Quotes_HeldTradeNotLastSaleNotConsolidated,
	"HELD_TRADE_LAST_SALE_BUT_NOT_ON_CONSOLIDATED":                  Financial_Quotes_HeldTradeLastSaleButNotConsolidated,
	"Held Trade Last Sale but not Consolidated":                     Financial_Quotes_HeldTradeLastSaleButNotConsolidated,
	"HELD_TRADE_LAST_SALE_AND_ON_CONSOLIDATED":                      Financial_Quotes_HeldTradeLastSaleAndConsolidated,
	"Held Trade Last Sale and Consolidated":                         Financial_Quotes_HeldTradeLastSaleAndConsolidated,
	"RETAIL_INTEREST_ON_BID":                                        Financial_Quotes_RetailInterestOnBid,
	"Retail Interest on Bid":                                        Financial_Quotes_RetailInterestOnBid,
	"RETAIL_INTEREST_ON_ASK":                                        Financial_Quotes_RetailInterestOnAsk,
	"Retail Interest on Ask":                                        Financial_Quotes_RetailInterestOnAsk,
	"RETAIL_INTEREST_ON_BID_AND_ASK":                                Financial_Quotes_RetailInterestOnBidAndAsk,
	"Retail Interest on Bid and Ask":                                Financial_Quotes_RetailInterestOnBidAndAsk,
	"FINRA_BBO_NO_CHANGE":                                           Financial_Quotes_FinraBBONoChange,
	"FINRA BBO: No Change":                                          Financial_Quotes_FinraBBONoChange,
	"FINRA_BBO_DOES_NOT_EXIST":                                      Financial_Quotes_FinraBBODoesNotExist,
	"FINRA BBO: Does not Exist":                                     Financial_Quotes_FinraBBODoesNotExist,
	"FINRA_BB_BO_EXECUTABLE":                                        Financial_Quotes_FinraBBBOExecutable,
	"FINRA BB / BO: Executable":                                     Financial_Quotes_FinraBBBOExecutable,
	"FINRA_BB_BELOW_LOWER_BAND":                                     Financial_Quotes_FinraBBBelowLowerBand,
	"FINRA BB: Below Lower Band":                                    Financial_Quotes_FinraBBBelowLowerBand,
	"FINRA_BO_ABOVE_UPPER_BAND":                                     Financial_Quotes_FinraBOAboveUpperBand,
	"FINRA BO: Above Upper Band":                                    Financial_Quotes_FinraBOAboveUpperBand,
	"FINRA_BB_BELOW_LOWER_BAND_BO_ABOVE_UPPER_BAND":                 Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand,
	"FINRA: BB Below Lower Band and BO Above Upper Band":            Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand,
	"CTA_NOT_DUE_TO_RELATED_SECURITY":                               Financial_Quotes_CTANotDueToRelatedSecurity,
	"CTA: Not Due to Related Security":                              Financial_Quotes_CTANotDueToRelatedSecurity,
	"CTA_DUE_TO_RELATED_SECURITY":                                   Financial_Quotes_CTADueToRelatedSecurity,
	"CTA: Due to Related Security":                                  Financial_Quotes_CTADueToRelatedSecurity,
	"CTA_NOT_IN_VIEW_OF_COMMON":                                     Financial_Quotes_CTANotInViewOfCommon,
	"CTA: Not in View of Common":                                    Financial_Quotes_CTANotInViewOfCommon,
	"CTA_IN_VIEW_OF_COMMON":                                         Financial_Quotes_CTAInViewOfCommon,
	"CTA: In View of Common":                                        Financial_Quotes_CTAInViewOfCommon,
	"CTA_PRICE_INDICATOR":                                           Financial_Quotes_CTAPriceIndicator,
	"CTA: Price Indicator":                                          Financial_Quotes_CTAPriceIndicator,
	"CTA_NEW_PRICE_INDICATOR":                                       Financial_Quotes_CTANewPriceIndicator,
	"CTA: New Price Indicator":                                      Financial_Quotes_CTANewPriceIndicator,
	"CTA_CORRECTED_PRICE_INDICATION":                                Financial_Quotes_CTACorrectedPriceIndication,
	"CTA: Corrected Price Indicator":                                Financial_Quotes_CTACorrectedPriceIndication,
	"CTA_CANCELLED_MARKET_IMBALANCE_PRICE_TRADING_RANGE_INDICATION": Financial_Quotes_CTACancelledMarketImbalance,
	"CTA: Cancelled Market Imbalance":                               Financial_Quotes_CTACancelledMarketImbalance,
}

// QuoteIndicatorMapping contains alternate names for the Financial.Quotes.Indicator enum
var QuoteIndicatorMapping = map[Financial_Quotes_Indicator]string{
	Financial_Quotes_NBBNBOExecutable:                            "NBB and/or NBO are Executable",
	Financial_Quotes_NBBBelowLowerBand:                           "NBB below Lower Band",
	Financial_Quotes_NBOAboveUpperBand:                           "NBO above Upper Band",
	Financial_Quotes_NBBBelowLowerBandAndNBOAboveUpperBand:       "NBB below Lower Band and NBO above Upper Band",
	Financial_Quotes_NBBEqualsUpperBand:                          "NBB equals Upper Band",
	Financial_Quotes_NBOEqualsLowerBand:                          "NBO equals Lower Band",
	Financial_Quotes_NBBEqualsUpperBandAndNBOAboveUpperBand:      "NBB equals Upper Band and NBO above Upper Band",
	Financial_Quotes_NBBBelowLowerBandAndNBOEqualsLowerBand:      "NBB below Lower Band and NBO equals Lower Band",
	Financial_Quotes_BidPriceAboveUpperLimitPriceBand:            "Bid Price above Upper Limit Price Band",
	Financial_Quotes_OfferPriceBelowLowerLimitPriceBand:          "Offer Price below Lower Limit Price Band",
	Financial_Quotes_BidAndOfferOutsidePriceBand:                 "Bid and Offer outside Price Band",
	Financial_Quotes_OpeningUpdate:                               "Opening Update",
	Financial_Quotes_IntraDayUpdate:                              "Intra-Day Update",
	Financial_Quotes_RestatedValue:                               "Restated Value",
	Financial_Quotes_SuspendedDuringTradingHalt:                  "Suspended during Trading Halt or Trading Pause",
	Financial_Quotes_ReOpeningUpdate:                             "Re-Opening Update",
	Financial_Quotes_OutsidePriceBandRuleHours:                   "Outside Price Band Rule Hours",
	Financial_Quotes_AuctionExtension:                            "Auction Extension (Auction Collar Message)",
	Financial_Quotes_LULDPriceBandInd:                            "LULD Price Band",
	Financial_Quotes_RepublishedLULDPriceBandInd:                 "Republished LULD Price Band",
	Financial_Quotes_NBBLimitStateEntered:                        "NBB Limit State Entered",
	Financial_Quotes_NBBLimitStateExited:                         "NBB Limit State Exited",
	Financial_Quotes_NBOLimitStateEntered:                        "NBO Limit State Entered",
	Financial_Quotes_NBOLimitStateExited:                         "NBO Limit State Exited",
	Financial_Quotes_NBBAndNBOLimitStateEntered:                  "NBB and NBO Limit State Entered",
	Financial_Quotes_NBBAndNBOLimitStateExited:                   "NBB and NBO Limit State Exited",
	Financial_Quotes_NBBLimitStateEnteredNBOLimitStateExited:     "NBB Limit State Entered and NBO Limit State Exited",
	Financial_Quotes_NBBLimitStateExitedNBOLimitStateEntered:     "NBB Limit State Exited and NBO Limit State Entered",
	Financial_Quotes_Deficient:                                   "Deficient - Below Listing Requirements",
	Financial_Quotes_Delinquent:                                  "Delinquent - Late Filing",
	Financial_Quotes_BankruptAndDeficient:                        "Bankrupt and Deficient",
	Financial_Quotes_BankruptAndDelinquent:                       "Bankrupt and Delinquent",
	Financial_Quotes_DeficientAndDelinquent:                      "Deficient and Delinquent",
	Financial_Quotes_DeficientDeliquentBankrupt:                  "Deficient, Delinquent, and Bankrupt",
	Financial_Quotes_CreationsSuspended:                          "Creations Suspended",
	Financial_Quotes_RedemptionsSuspended:                        "Redemptions Suspended",
	Financial_Quotes_CreationsRedemptionsSuspended:               "Creations and/or Redemptions Suspended",
	Financial_Quotes_NormalTrading:                               "Normal Trading",
	Financial_Quotes_OpeningDelay:                                "Opening Delay",
	Financial_Quotes_TradingHalt:                                 "Trading Halt",
	Financial_Quotes_TradingResume:                               "Resume",
	Financial_Quotes_NoOpenNoResume:                              "No Open / No Resume",
	Financial_Quotes_PriceIndication:                             "Price Indication",
	Financial_Quotes_TradingRangeIndication:                      "Trading Range Indication",
	Financial_Quotes_MarketImbalanceBuy:                          "Market Imbalance Buy",
	Financial_Quotes_MarketImbalanceSell:                         "Market Imbalance Sell",
	Financial_Quotes_MarketOnCloseImbalanceBuy:                   "Market On-Close Imbalance Buy",
	Financial_Quotes_MarketOnCloseImbalanceSell:                  "Market On Close Imbalance Sell",
	Financial_Quotes_NoMarketImbalance:                           "No Market Imbalance",
	Financial_Quotes_NoMarketOnCloseImbalance:                    "No Market, On-Close Imbalance",
	Financial_Quotes_ShortSaleRestriction:                        "Short Sale Restriction",
	Financial_Quotes_LimitUpLimitDown:                            "Limit Up-Limit Down",
	Financial_Quotes_QuotationResumption:                         "Quotation Resumption",
	Financial_Quotes_TradingResumption:                           "Trading Resumption",
	Financial_Quotes_VolatilityTradingPause:                      "Volatility Trading Pause",
	Financial_Quotes_HaltNewsPending:                             "Halt: News Pending",
	Financial_Quotes_UpdateNewsDissemination:                     "Update: News Dissemination",
	Financial_Quotes_HaltSingleStockTradingPause:                 "Halt: Single Stock Trading Pause in Affect",
	Financial_Quotes_HaltRegulatoryExtraordinaryMarketActivity:   "Halt: Regulatory Extraordinary Market Activity",
	Financial_Quotes_HaltETF:                                     "Halt: ETF",
	Financial_Quotes_HaltInformationRequested:                    "Halt: Information Requested",
	Financial_Quotes_HaltExchangeNonCompliance:                   "Halt: Exchange Non-Compliance",
	Financial_Quotes_HaltFilingsNotCurrent:                       "Halt: Filings Not Current",
	Financial_Quotes_HaltSECTradingSuspension:                    "Halt: SEC Trading Suspension",
	Financial_Quotes_HaltRegulatoryConcern:                       "Halt: Regulatory Concern",
	Financial_Quotes_HaltMarketOperations:                        "Halt: Market Operations",
	Financial_Quotes_IPOSecurityNotYetTrading:                    "IPO Security: Not Yet Trading",
	Financial_Quotes_HaltCorporateAction:                         "Halt: Corporate Action",
	Financial_Quotes_QuotationNotAvailable:                       "Quotation Not Available",
	Financial_Quotes_HaltVolatilityTradingPause:                  "Halt: Volatility Trading Pause",
	Financial_Quotes_HaltVolatilityTradingPauseStraddleCondition: "Halt: Volatility Trading Pause - Straddle Condition",
	Financial_Quotes_UpdateNewsAndResumptionTimes:                "Update: News and Resumption Times",
	Financial_Quotes_HaltSingleStockTradingPauseQuotesOnly:       "Halt: Single Stock Trading Pause - Quotes Only",
	Financial_Quotes_ResumeQualificationIssuesReviewedResolved:   "Resume: Qualification Issues Reviewed / Resolved",
	Financial_Quotes_ResumeFilingRequirementsSatisfiedResolved:   "Resume: Filing Requirements Satisfied / Resolved",
	Financial_Quotes_ResumeNewsNotForthcoming:                    "Resume: News Not Forthcoming",
	Financial_Quotes_ResumeQualificationsMaintRequirementsMet:    "Resume: Qualifications - Maintenance Requirements Met",
	Financial_Quotes_ResumeQualificationsFilingsMet:              "Resume: Qualifications - Filings Met",
	Financial_Quotes_ResumeRegulatoryAuth:                        "Resume: Regulatory Auth",
	Financial_Quotes_NewIssueAvailable:                           "New Issue Available",
	Financial_Quotes_IssueAvailable:                              "Issue Available",
	Financial_Quotes_MWCBCarryFromPreviousDay:                    "MWCB - Carry from Previous Day",
	Financial_Quotes_MWCBResume:                                  "MWCB - Resume",
	Financial_Quotes_IPOSecurityReleasedForQuotation:             "IPO Security: Released for Quotation",
	Financial_Quotes_IPOSecurityPositioningWindowExtension:       "IPO Security: Positioning Window Extension",
	Financial_Quotes_MWCBLevel1:                                  "MWCB - Level 1",
	Financial_Quotes_MWCBLevel2:                                  "MWCB - Level 2",
	Financial_Quotes_MWCBLevel3:                                  "MWCB - Level 3",
	Financial_Quotes_HaltSubPennyTrading:                         "Halt: Sub-Penny Trading",
	Financial_Quotes_OrderImbalanceInd:                           "Order Imbalance",
	Financial_Quotes_LULDTradingPaused:                           "LULD Trading Paused",
	Financial_Quotes_NONE:                                        "Security Status: None",
	Financial_Quotes_ShortSalesRestrictionActivated:              "Short Sales Restriction Activated",
	Financial_Quotes_ShortSalesRestrictionContinued:              "Short Sales Restriction Continued",
	Financial_Quotes_ShortSalesRestrictionDeactivated:            "Short Sales Restriction Deactivated",
	Financial_Quotes_ShortSalesRestrictionInEffect:               "Short Sales Restriction in Effect",
	Financial_Quotes_ShortSalesRestrictionMax:                    "Short Sales Restriction Max",
	Financial_Quotes_NBBONoChange:                                "NBBO: No Change",
	Financial_Quotes_NBBOQuoteIsNBBO:                             "NBBO: Quote is NBBO",
	Financial_Quotes_NBBONoBBNoBO:                                "NBBO: No BB, No BO",
	Financial_Quotes_NBBOBBBOShortAppendage:                      "NBBO: BB / BO Short Appendage",
	Financial_Quotes_NBBOBBBOLongAppendage:                       "NBBO: BB / BO Long Appendage",
	Financial_Quotes_HeldTradeNotLastSaleNotConsolidated:         "Held Trade not Last Sale, not Consolidated",
	Financial_Quotes_HeldTradeLastSaleButNotConsolidated:         "Held Trade Last Sale but not Consolidated",
	Financial_Quotes_HeldTradeLastSaleAndConsolidated:            "Held Trade Last Sale and Consolidated",
	Financial_Quotes_RetailInterestOnBid:                         "Retail Interest on Bid",
	Financial_Quotes_RetailInterestOnAsk:                         "Retail Interest on Ask",
	Financial_Quotes_RetailInterestOnBidAndAsk:                   "Retail Interest on Bid and Ask",
	Financial_Quotes_FinraBBONoChange:                            "FINRA BBO: No Change",
	Financial_Quotes_FinraBBODoesNotExist:                        "FINRA BBO: Does not Exist",
	Financial_Quotes_FinraBBBOExecutable:                         "FINRA BB / BO: Executable",
	Financial_Quotes_FinraBBBelowLowerBand:                       "FINRA BB: Below Lower Band",
	Financial_Quotes_FinraBOAboveUpperBand:                       "FINRA BO: Above Upper Band",
	Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand:      "FINRA: BB Below Lower Band and BO Above Upper Band",
	Financial_Quotes_CTANotDueToRelatedSecurity:                  "CTA: Not Due to Related Security",
	Financial_Quotes_CTADueToRelatedSecurity:                     "CTA: Due to Related Security",
	Financial_Quotes_CTANotInViewOfCommon:                        "CTA: Not in View of Common",
	Financial_Quotes_CTAInViewOfCommon:                           "CTA: In View of Common",
	Financial_Quotes_CTAPriceIndicator:                           "CTA: Price Indicator",
	Financial_Quotes_CTANewPriceIndicator:                        "CTA: New Price Indicator",
	Financial_Quotes_CTACorrectedPriceIndication:                 "CTA: Corrected Price Indicator",
	Financial_Quotes_CTACancelledMarketImbalance:                 "CTA: Cancelled Market Imbalance",
}

// TradeConditionAlternates contains alternative values for the Financial.Trades.Condition enum
var TradeConditionAlternates = map[string]Financial_Trades_Condition{
	"UTP:A":                  Financial_Trades_Acquisition,
	"UTP:W":                  Financial_Trades_AveragePriceTrade,
	"UTP:B":                  Financial_Trades_BunchedTrade,
	"UTP:G":                  Financial_Trades_BunchedSoldTrade,
	"UTP:C":                  Financial_Trades_CashSale,
	"UTP:6":                  Financial_Trades_ClosingPrints,
	"UTP:X":                  Financial_Trades_CrossTrade,
	"UTP:4":                  Financial_Trades_DerivativelyPriced,
	"UTP:D":                  Financial_Trades_Distribution,
	"UTP:T":                  Financial_Trades_FormT,
	"UTP:U":                  Financial_Trades_ExtendedTradingHours,
	"UTP:F":                  Financial_Trades_IntermarketSweep,
	"UTP:M":                  Financial_Trades_MarketCenterOfficialClose,
	"UTP:Q":                  Financial_Trades_MarketCenterOfficialOpen,
	"UTP:N":                  Financial_Trades_NextDay,
	"UTP:H":                  Financial_Trades_PriceVariationTrade,
	"UTP:P":                  Financial_Trades_PriorReferencePrice,
	"UTP:K":                  Financial_Trades_Rule155Trade,
	"UTP:O":                  Financial_Trades_OpeningPrints,
	"UTP:1":                  Financial_Trades_StoppedStock,
	"UTP:R":                  Financial_Trades_Seller,
	"UTP:5":                  Financial_Trades_ReOpeningPrints,
	"UTP:L":                  Financial_Trades_SoldLast,
	"UTP:2":                  Financial_Trades_SoldLastAndStoppedStock,
	"UTP:Z":                  Financial_Trades_SoldOut,
	"UTP:3":                  Financial_Trades_SoldOutOfSequence,
	"UTP:S":                  Financial_Trades_SplitTrade,
	"UTP:V":                  Financial_Trades_StockOption,
	"UTP:Y":                  Financial_Trades_YellowFlagRegularTrade,
	"UTP:I":                  Financial_Trades_OddLotTrade,
	"UTP:9":                  Financial_Trades_CorrectedConsolidatedClose,
	"UTP:7":                  Financial_Trades_QualifiedContingentTrade,
	"CTA:B":                  Financial_Trades_AveragePriceTrade,
	"CTA:E":                  Financial_Trades_AutomaticExecution,
	"CTA:I":                  Financial_Trades_CAPElection,
	"CTA:C":                  Financial_Trades_CashSale,
	"CTA:X":                  Financial_Trades_CrossTrade,
	"CTA:4":                  Financial_Trades_DerivativelyPriced,
	"CTA:T":                  Financial_Trades_FormT,
	"CTA:U":                  Financial_Trades_ExtendedTradingHours,
	"CTA:F":                  Financial_Trades_IntermarketSweep,
	"CTA:M":                  Financial_Trades_MarketCenterOfficialClose,
	"CTA:Q":                  Financial_Trades_MarketCenterOfficialOpen,
	"CTA:O":                  Financial_Trades_MarketCenterOpeningTrade,
	"CTA:S":                  Financial_Trades_MarketCenterReopeningTrade,
	"CTA:6":                  Financial_Trades_MarketCenterClosingTrade,
	"CTA:N":                  Financial_Trades_NextDay,
	"CTA:H":                  Financial_Trades_PriceVariationTrade,
	"CTA:P":                  Financial_Trades_PriorReferencePrice,
	"CTA:K":                  Financial_Trades_Rule155Trade,
	"CTA:R":                  Financial_Trades_Seller,
	"CTA:L":                  Financial_Trades_SoldLast,
	"CTA:Z":                  Financial_Trades_SoldOut,
	"CTA:9":                  Financial_Trades_CorrectedConsolidatedClose,
	"CTA:1":                  Financial_Trades_TradeThruExempt,
	"CTA:V":                  Financial_Trades_ContingentTrade,
	"CTA:7":                  Financial_Trades_QualifiedContingentTrade,
	"CTA:G":                  Financial_Trades_OpeningReopeningTradeDetail,
	"CTA:A":                  Financial_Trades_ShortSaleRestrictionActivated,
	"CTA:D":                  Financial_Trades_ShortSaleRestrictionDeactivated,
	"CTA:2":                  Financial_Trades_FinancialStatusDeficient,
	"CTA:3":                  Financial_Trades_FinancialStatusDelinquent,
	"CTA:5":                  Financial_Trades_FinancialStatusBankruptAndDelinquent,
	"CTA:8":                  Financial_Trades_FinancialStatusCreationsSuspended,
	"FINRA_TDDS:W":           Financial_Trades_AveragePriceTrade,
	"FINRA_TDDS:C":           Financial_Trades_CashSale,
	"FINRA_TDDS:T":           Financial_Trades_FormT,
	"FINRA_TDDS:U":           Financial_Trades_ExtendedTradingHours,
	"FINRA_TDDS:N":           Financial_Trades_NextDay,
	"FINRA_TDDS:P":           Financial_Trades_PriorReferencePrice,
	"FINRA_TDDS:R":           Financial_Trades_Seller,
	"FINRA_TDDS:Z":           Financial_Trades_SoldOut,
	"FINRA_TDDS:I":           Financial_Trades_OddLotTrade,
	"OPRA:A":                 Financial_Trades_Canceled,
	"OPRA:B":                 Financial_Trades_LateAndOutOfSequence,
	"OPRA:C":                 Financial_Trades_LastAndCanceled,
	"OPRA:D":                 Financial_Trades_Late,
	"OPRA:E":                 Financial_Trades_OpeningTradeAndCanceled,
	"OPRA:F":                 Financial_Trades_OpeningTradeLateAndOutOfSequence,
	"OPRA:G":                 Financial_Trades_OnlyTradeAndCanceled,
	"OPRA:H":                 Financial_Trades_OpeningTradeAndLate,
	"OPRA:I":                 Financial_Trades_AutomaticExecutionOption,
	"OPRA:J":                 Financial_Trades_ReopeningTrade,
	"OPRA:S":                 Financial_Trades_IntermarketSweepOrder,
	"OPRA:a":                 Financial_Trades_SingleLegAuctionNonISO,
	"OPRA:b":                 Financial_Trades_SingleLegAuctionISO,
	"OPRA:c":                 Financial_Trades_SingleLegCrossNonISO,
	"OPRA:d":                 Financial_Trades_SingleLegCrossISO,
	"OPRA:e":                 Financial_Trades_SingleLegFloorTrade,
	"OPRA:f":                 Financial_Trades_MultiLegAutoElectronicTrade,
	"OPRA:g":                 Financial_Trades_MultiLegAuction,
	"OPRA:h":                 Financial_Trades_MultiLegCross,
	"OPRA:i":                 Financial_Trades_MultiLegFloorTrade,
	"OPRA:j":                 Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg,
	"OPRA:k":                 Financial_Trades_StockOptionsAuction,
	"OPRA:l":                 Financial_Trades_MultiLegAuctionAgainstSingleLeg,
	"OPRA:m":                 Financial_Trades_MultiLegFloorTradeAgainstSingleLeg,
	"OPRA:n":                 Financial_Trades_StockOptionsAutoElectronicTrade,
	"OPRA:o":                 Financial_Trades_StockOptionsCross,
	"OPRA:p":                 Financial_Trades_StockOptionsFloorTrade,
	"OPRA:q":                 Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg,
	"OPRA:r":                 Financial_Trades_StockOptionsAuctionAgainstSingleLeg,
	"OPRA:s":                 Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg,
	"OPRA:t":                 Financial_Trades_MultiLegFloorTradeOfProprietaryProducts,
	"OPRA:u":                 Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts,
	"OPRA:v":                 Financial_Trades_ExtendedHoursTrade,
	"CANC":                   Financial_Trades_Canceled,
	"OSEQ":                   Financial_Trades_LateAndOutOfSequence,
	"CNCL":                   Financial_Trades_LastAndCanceled,
	"LATE":                   Financial_Trades_Late,
	"CNCO":                   Financial_Trades_OpeningTradeAndCanceled,
	"OPEN":                   Financial_Trades_OpeningTradeLateAndOutOfSequence,
	"CNOL":                   Financial_Trades_OnlyTradeAndCanceled,
	"OPNL":                   Financial_Trades_OpeningTradeAndLate,
	"AUTO":                   Financial_Trades_AutomaticExecutionOption,
	"REOP":                   Financial_Trades_ReopeningTrade,
	"ISOI":                   Financial_Trades_IntermarketSweepOrder,
	"SLAN":                   Financial_Trades_SingleLegAuctionNonISO,
	"SLAI":                   Financial_Trades_SingleLegAuctionISO,
	"SLCN":                   Financial_Trades_SingleLegCrossNonISO,
	"SLCI":                   Financial_Trades_SingleLegCrossISO,
	"SLFT":                   Financial_Trades_SingleLegFloorTrade,
	"MLET":                   Financial_Trades_MultiLegAutoElectronicTrade,
	"MLAT":                   Financial_Trades_MultiLegAuction,
	"MLCT":                   Financial_Trades_MultiLegCross,
	"MLFT":                   Financial_Trades_MultiLegFloorTrade,
	"MESL":                   Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg,
	"TLAT":                   Financial_Trades_StockOptionsAuction,
	"MASL":                   Financial_Trades_MultiLegAuctionAgainstSingleLeg,
	"MFSL":                   Financial_Trades_MultiLegFloorTradeAgainstSingleLeg,
	"TLET":                   Financial_Trades_StockOptionsAutoElectronicTrade,
	"TLCT":                   Financial_Trades_StockOptionsCross,
	"TLFT":                   Financial_Trades_StockOptionsFloorTrade,
	"TESL":                   Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg,
	"TASL":                   Financial_Trades_StockOptionsAuctionAgainstSingleLeg,
	"TFSL":                   Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg,
	"CBMO":                   Financial_Trades_MultiLegFloorTradeOfProprietaryProducts,
	"MCTP":                   Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts,
	"EXHT":                   Financial_Trades_ExtendedHoursTrade,
	"Regular Sale":           Financial_Trades_RegularSale,
	"Average Price Trade":    Financial_Trades_AveragePriceTrade,
	"Automatic Execution":    Financial_Trades_AutomaticExecution,
	"Bunched Trade":          Financial_Trades_BunchedTrade,
	"Bunched Sold Trade":     Financial_Trades_BunchedSoldTrade,
	"CAP Election":           Financial_Trades_CAPElection,
	"Cash Sale":              Financial_Trades_CashSale,
	"Closing Prints":         Financial_Trades_ClosingPrints,
	"Cross Trade":            Financial_Trades_CrossTrade,
	"Derivatively Priced":    Financial_Trades_DerivativelyPriced,
	"Form T":                 Financial_Trades_FormT,
	"Extended Trading Hours": Financial_Trades_ExtendedTradingHours,
	"Sold Out of Sequence":   Financial_Trades_ExtendedTradingHours,
	"Extended Trading Hours (Sold Out of Sequence)":           Financial_Trades_ExtendedTradingHours,
	"Intermarket Sweep":                                       Financial_Trades_IntermarketSweep,
	"Market Center Official Close":                            Financial_Trades_MarketCenterOfficialClose,
	"Market Center Official Open":                             Financial_Trades_MarketCenterOfficialOpen,
	"Market Center Opening Trade":                             Financial_Trades_MarketCenterOpeningTrade,
	"Market Center Reopening Trade":                           Financial_Trades_MarketCenterReopeningTrade,
	"Market Center Closing Trade":                             Financial_Trades_MarketCenterClosingTrade,
	"Next Day":                                                Financial_Trades_NextDay,
	"Price Variation Trade":                                   Financial_Trades_PriceVariationTrade,
	"Prior Reference Price":                                   Financial_Trades_PriorReferencePrice,
	"Rule 155 Trade (AMEX)":                                   Financial_Trades_Rule155Trade,
	"Rule 127 NYSE":                                           Financial_Trades_Rule127NYSE,
	"Opening Prints":                                          Financial_Trades_OpeningPrints,
	"Stopped Stock (Regular Trade)":                           Financial_Trades_StoppedStock,
	"Re-Opening Prints":                                       Financial_Trades_ReOpeningPrints,
	"Sold Last":                                               Financial_Trades_SoldLast,
	"Sold Last and Stopped Stock":                             Financial_Trades_SoldLastAndStoppedStock,
	"Sold Out":                                                Financial_Trades_SoldOut,
	"Sold (Out of Sequence)":                                  Financial_Trades_SoldOutOfSequence,
	"Split Trade":                                             Financial_Trades_SplitTrade,
	"Stock Option":                                            Financial_Trades_StockOption,
	"Yellow Flag Regular Trade":                               Financial_Trades_YellowFlagRegularTrade,
	"Odd Lot Trade":                                           Financial_Trades_OddLotTrade,
	"Corrected Consolidated Close":                            Financial_Trades_CorrectedConsolidatedClose,
	"Trade Thru Exempt":                                       Financial_Trades_TradeThruExempt,
	"Non-Eligible":                                            Financial_Trades_NonEligible,
	"Non-Eligible Extended":                                   Financial_Trades_NonEligibleExtended,
	"As of":                                                   Financial_Trades_AsOf,
	"As of Correction":                                        Financial_Trades_AsOfCorrection,
	"As of Cancel":                                            Financial_Trades_AsOfCancel,
	"Contingent Trade":                                        Financial_Trades_ContingentTrade,
	"Qualified Contingent Trade (QCT)":                        Financial_Trades_QualifiedContingentTrade,
	"OPENING_REOPENING_TRADE_DETAIL":                          Financial_Trades_OpeningReopeningTradeDetail,
	"Opening / Reopening Trade Detail":                        Financial_Trades_OpeningReopeningTradeDetail,
	"Short Sale Restriction Activated":                        Financial_Trades_ShortSaleRestrictionActivated,
	"Short Sale Restriction Continued":                        Financial_Trades_ShortSaleRestrictionContinued,
	"Short Sale Restriction Deactivated":                      Financial_Trades_ShortSaleRestrictionDeactivated,
	"Short Sale Restriction in Effect":                        Financial_Trades_ShortSaleRestrictionInEffect,
	"Financial Status: Bankrupt":                              Financial_Trades_FinancialStatusBankrupt,
	"Financial Status: Deficient":                             Financial_Trades_FinancialStatusDeficient,
	"Financial Status: Delinquent":                            Financial_Trades_FinancialStatusDelinquent,
	"Financial Status: Bankrupt and Deficient":                Financial_Trades_FinancialStatusBankruptAndDeficient,
	"Financial Status: Bankrupt and Delinquent":               Financial_Trades_FinancialStatusBankruptAndDelinquent,
	"Financial Status: Deficient and Delinquent":              Financial_Trades_FinancialStatusDeficientAndDelinquent,
	"Financial Status: Deficient, Delinquent, Bankrupt":       Financial_Trades_FinancialStatusDeficientDelinquentBankrupt,
	"Financial Status: Liquidation":                           Financial_Trades_FinancialStatusLiquidation,
	"Financial Status: Creations Suspended":                   Financial_Trades_FinancialStatusCreationsSuspended,
	"Financial Status: Redemptions Suspended":                 Financial_Trades_FinancialStatusRedemptionsSuspended,
	"Late and Out of Sequence":                                Financial_Trades_LateAndOutOfSequence,
	"Last and Canceled":                                       Financial_Trades_LastAndCanceled,
	"Opening Trade and Canceled":                              Financial_Trades_OpeningTradeAndCanceled,
	"Opening Trade, Late and Out of Sequence":                 Financial_Trades_OpeningTradeLateAndOutOfSequence,
	"Only Trade and Canceled":                                 Financial_Trades_OnlyTradeAndCanceled,
	"Opening Trade and Late":                                  Financial_Trades_OpeningTradeAndLate,
	"Automatic Execution Option":                              Financial_Trades_AutomaticExecutionOption,
	"Reopening Trade":                                         Financial_Trades_ReopeningTrade,
	"Intermarket Sweep Order":                                 Financial_Trades_IntermarketSweepOrder,
	"Single-Leg Auction, Non-ISO":                             Financial_Trades_SingleLegAuctionNonISO,
	"Single-Leg Auction, ISO":                                 Financial_Trades_SingleLegAuctionISO,
	"Single-Leg Cross, Non-ISO":                               Financial_Trades_SingleLegCrossNonISO,
	"Single-Leg Cross, ISO":                                   Financial_Trades_SingleLegCrossISO,
	"Single-Leg Floor Trade":                                  Financial_Trades_SingleLegFloorTrade,
	"Multi-Leg, Auto-Electronic Trade":                        Financial_Trades_MultiLegAutoElectronicTrade,
	"Multi-Leg Auction":                                       Financial_Trades_MultiLegAuction,
	"Multi-Leg Cross":                                         Financial_Trades_MultiLegCross,
	"Multi-Leg Floor Trade":                                   Financial_Trades_MultiLegFloorTrade,
	"Multi-Leg, Auto-Electronic Trade against Single-Leg":     Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg,
	"Stock Options Auction":                                   Financial_Trades_StockOptionsAuction,
	"Multi-Leg Auction against Single-Leg":                    Financial_Trades_MultiLegAuctionAgainstSingleLeg,
	"Multi-Leg Floor Trade against Single-Leg":                Financial_Trades_MultiLegFloorTradeAgainstSingleLeg,
	"Stock Options, Auto-Electronic Trade":                    Financial_Trades_StockOptionsAutoElectronicTrade,
	"Stock Options Cross":                                     Financial_Trades_StockOptionsCross,
	"Stock Options Floor Trade":                               Financial_Trades_StockOptionsFloorTrade,
	"Stock Options, Auto-Electronic Trade against Single-Leg": Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg,
	"Stock Options, Auction against Single-Leg":               Financial_Trades_StockOptionsAuctionAgainstSingleLeg,
	"Stock Options, Floor Trade against Single-Leg":           Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg,
	"Multi-Leg Floor Trade of Proprietary Products":           Financial_Trades_MultiLegFloorTradeOfProprietaryProducts,
	"Multilateral Compression Trade of Proprietary Products":  Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts,
	"Extended Hours Trade":                                    Financial_Trades_ExtendedHoursTrade,
}

// TradeConditionMapping contains alternate names for the Financial.Trades.Condition enum
var TradeConditionMapping = map[Financial_Trades_Condition]string{
	Financial_Trades_RegularSale:                                       "Regular Sale",
	Financial_Trades_AveragePriceTrade:                                 "Average Price Trade",
	Financial_Trades_AutomaticExecution:                                "Automatic Execution",
	Financial_Trades_BunchedTrade:                                      "Bunched Trade",
	Financial_Trades_BunchedSoldTrade:                                  "Bunched Sold Trade",
	Financial_Trades_CAPElection:                                       "CAP Election",
	Financial_Trades_CashSale:                                          "Cash Sale",
	Financial_Trades_ClosingPrints:                                     "Closing Prints",
	Financial_Trades_CrossTrade:                                        "Cross Trade",
	Financial_Trades_DerivativelyPriced:                                "Derivatively Priced",
	Financial_Trades_FormT:                                             "Form T",
	Financial_Trades_ExtendedTradingHours:                              "Extended Trading Hours (Sold Out of Sequence)",
	Financial_Trades_IntermarketSweep:                                  "Intermarket Sweep",
	Financial_Trades_MarketCenterOfficialClose:                         "Market Center Official Close",
	Financial_Trades_MarketCenterOfficialOpen:                          "Market Center Official Open",
	Financial_Trades_MarketCenterOpeningTrade:                          "Market Center Opening Trade",
	Financial_Trades_MarketCenterReopeningTrade:                        "Market Center Reopening Trade",
	Financial_Trades_MarketCenterClosingTrade:                          "Market Center Closing Trade",
	Financial_Trades_NextDay:                                           "Next Day",
	Financial_Trades_PriceVariationTrade:                               "Price Variation Trade",
	Financial_Trades_PriorReferencePrice:                               "Prior Reference Price",
	Financial_Trades_Rule155Trade:                                      "Rule 155 Trade (AMEX)",
	Financial_Trades_Rule127NYSE:                                       "Rule 127 NYSE",
	Financial_Trades_OpeningPrints:                                     "Opening Prints",
	Financial_Trades_StoppedStock:                                      "Stopped Stock (Regular Trade)",
	Financial_Trades_ReOpeningPrints:                                   "Re-Opening Prints",
	Financial_Trades_SoldLast:                                          "Sold Last",
	Financial_Trades_SoldLastAndStoppedStock:                           "Sold Last and Stopped Stock",
	Financial_Trades_SoldOut:                                           "Sold Out",
	Financial_Trades_SoldOutOfSequence:                                 "Sold (Out of Sequence)",
	Financial_Trades_SplitTrade:                                        "Split Trade",
	Financial_Trades_StockOption:                                       "Stock Option",
	Financial_Trades_YellowFlagRegularTrade:                            "Yellow Flag Regular Trade",
	Financial_Trades_OddLotTrade:                                       "Odd Lot Trade",
	Financial_Trades_CorrectedConsolidatedClose:                        "Corrected Consolidated Close",
	Financial_Trades_TradeThruExempt:                                   "Trade Thru Exempt",
	Financial_Trades_NonEligible:                                       "Non-Eligible",
	Financial_Trades_NonEligibleExtended:                               "Non-Eligible Extended",
	Financial_Trades_AsOf:                                              "As of",
	Financial_Trades_AsOfCorrection:                                    "As of Correction",
	Financial_Trades_AsOfCancel:                                        "As of Cancel",
	Financial_Trades_ContingentTrade:                                   "Contingent Trade",
	Financial_Trades_QualifiedContingentTrade:                          "Qualified Contingent Trade (QCT)",
	Financial_Trades_OpeningReopeningTradeDetail:                       "Opening / Reopening Trade Detail",
	Financial_Trades_ShortSaleRestrictionActivated:                     "Short Sale Restriction Activated",
	Financial_Trades_ShortSaleRestrictionContinued:                     "Short Sale Restriction Continued",
	Financial_Trades_ShortSaleRestrictionDeactivated:                   "Short Sale Restriction Deactivated",
	Financial_Trades_ShortSaleRestrictionInEffect:                      "Short Sale Restriction in Effect",
	Financial_Trades_FinancialStatusBankrupt:                           "Financial Status: Bankrupt",
	Financial_Trades_FinancialStatusDeficient:                          "Financial Status: Deficient",
	Financial_Trades_FinancialStatusDelinquent:                         "Financial Status: Delinquent",
	Financial_Trades_FinancialStatusBankruptAndDeficient:               "Financial Status: Bankrupt and Deficient",
	Financial_Trades_FinancialStatusBankruptAndDelinquent:              "Financial Status: Bankrupt and Delinquent",
	Financial_Trades_FinancialStatusDeficientAndDelinquent:             "Financial Status: Deficient and Delinquent",
	Financial_Trades_FinancialStatusDeficientDelinquentBankrupt:        "Financial Status: Deficient, Delinquent, Bankrupt",
	Financial_Trades_FinancialStatusLiquidation:                        "Financial Status: Liquidation",
	Financial_Trades_FinancialStatusCreationsSuspended:                 "Financial Status: Creations Suspended",
	Financial_Trades_FinancialStatusRedemptionsSuspended:               "Financial Status: Redemptions Suspended",
	Financial_Trades_LateAndOutOfSequence:                              "Late and Out of Sequence",
	Financial_Trades_LastAndCanceled:                                   "Last and Canceled",
	Financial_Trades_OpeningTradeAndCanceled:                           "Opening Trade and Canceled",
	Financial_Trades_OpeningTradeLateAndOutOfSequence:                  "Opening Trade, Late and Out of Sequence",
	Financial_Trades_OnlyTradeAndCanceled:                              "Only Trade and Canceled",
	Financial_Trades_OpeningTradeAndLate:                               "Opening Trade and Late",
	Financial_Trades_AutomaticExecutionOption:                          "Automatic Execution Option",
	Financial_Trades_ReopeningTrade:                                    "Reopening Trade",
	Financial_Trades_IntermarketSweepOrder:                             "Intermarket Sweep Order",
	Financial_Trades_SingleLegAuctionNonISO:                            "Single-Leg Auction, Non-ISO",
	Financial_Trades_SingleLegAuctionISO:                               "Single-Leg Auction, ISO",
	Financial_Trades_SingleLegCrossNonISO:                              "Single-Leg Cross, Non-ISO",
	Financial_Trades_SingleLegCrossISO:                                 "Single-Leg Cross, ISO",
	Financial_Trades_SingleLegFloorTrade:                               "Single-Leg Floor Trade",
	Financial_Trades_MultiLegAutoElectronicTrade:                       "Multi-Leg, Auto-Electronic Trade",
	Financial_Trades_MultiLegAuction:                                   "Multi-Leg Auction",
	Financial_Trades_MultiLegCross:                                     "Multi-Leg Cross",
	Financial_Trades_MultiLegFloorTrade:                                "Multi-Leg Floor Trade",
	Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg:       "Multi-Leg, Auto-Electronic Trade against Single-Leg",
	Financial_Trades_StockOptionsAuction:                               "Stock Options Auction",
	Financial_Trades_MultiLegAuctionAgainstSingleLeg:                   "Multi-Leg Auction against Single-Leg",
	Financial_Trades_MultiLegFloorTradeAgainstSingleLeg:                "Multi-Leg Floor Trade against Single-Leg",
	Financial_Trades_StockOptionsAutoElectronicTrade:                   "Stock Options, Auto-Electronic Trade",
	Financial_Trades_StockOptionsCross:                                 "Stock Options Cross",
	Financial_Trades_StockOptionsFloorTrade:                            "Stock Options Floor Trade",
	Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg:   "Stock Options, Auto-Electronic Trade against Single-Leg",
	Financial_Trades_StockOptionsAuctionAgainstSingleLeg:               "Stock Options, Auction against Single-Leg",
	Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg:            "Stock Options, Floor Trade against Single-Leg",
	Financial_Trades_MultiLegFloorTradeOfProprietaryProducts:           "Multi-Leg Floor Trade of Proprietary Products",
	Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts: "Multilateral Compression Trade of Proprietary Products",
	Financial_Trades_ExtendedHoursTrade:                                "Extended Hours Trade",
}

// TradeCorrectionAlternates contains alternative valus for the Financial.Trades.CorrectionCode enum
var TradeCorrectionAlternates = map[string]Financial_Trades_CorrectionCode{
	"Not Corrected":     Financial_Trades_NotCorrected,
	"Late, Corrected":   Financial_Trades_LateCorrected,
	"Cancelled":         Financial_Trades_Cancel,
	"Cancel Record":     Financial_Trades_CancelRecord,
	"Error Record":      Financial_Trades_ErrorRecord,
	"Correction Record": Financial_Trades_CorrectionRecord,
	"00":                Financial_Trades_NotCorrected,
	"01":                Financial_Trades_LateCorrected,
	"07":                Financial_Trades_Erroneous,
	"08":                Financial_Trades_Cancel,
}

// TradeCorrectionMapping contains alternate names for the Financial.Trades.CorrectionCode enum
var TradeCorrectionMapping = map[Financial_Trades_CorrectionCode]string{
	Financial_Trades_NotCorrected:     "Not Corrected",
	Financial_Trades_LateCorrected:    "Late, Corrected",
	Financial_Trades_Cancel:           "Cancelled",
	Financial_Trades_CancelRecord:     "Cancel Record",
	Financial_Trades_ErrorRecord:      "Error Record",
	Financial_Trades_CorrectionRecord: "Correction Record",
}

// MarhsalJSON converts a Decimal to JSON
func (d *Decimal) MarshalJSON() ([]byte, error) {
	return []byte(d.ToString()), nil
}

// MarshalCSV converts a Decimal to a CSV format
func (d *Decimal) MarshalCSV() (string, error) {
	return d.ToString(), nil
}

// Marshaler converts a Decimal to a DynamoDB attribute value
func (d *Decimal) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberN{
		Value: d.ToString(),
	}, nil
}

// Value converts a Decimal to an SQL value
func (d *Decimal) Value() (driver.Value, error) {
	return driver.Value(d.ToString()), nil
}

// UnmarshalJSON converts JSON data into a Decimal
func (d *Decimal) UnmarshalJSON(data []byte) error {

	// Check if the value is nil; if this is the case then return nil
	if data == nil {
		return nil
	}

	// Otherwise, convert the data from a string into a timestamp
	return d.FromString(string(data))
}

// UnmarshalCSV converts a CSV column into a Decimal
func (d *Decimal) UnmarshalCSV(raw string) error {
	return d.FromString(raw)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Decimal
func (d *Decimal) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return d.FromString(string(casted.Value))
	case *types.AttributeValueMemberN:
		return d.FromString(casted.Value)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return d.FromString(casted.Value)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Decimal", value)
	}
}

// Scan converts an SQL value into a Decimal
func (d *Decimal) Scan(value interface{}) error {

	// Check if the value is nil; if this is the case then return nil
	if value == nil {
		return nil
	}

	// Based on the type of the value we're working with, we'll convert the decimal from its implied
	// type to a Decimal; if this fails or the type isn't one we recognized then we'll return an error
	switch casted := value.(type) {
	case []byte:
		return d.FromString(string(casted))
	case float64:
		*d = *NewFromDecimal(decimal.NewFromFloat(casted))
	case int64:
		*d = *NewFromDecimal(decimal.NewFromInt(casted))
	case string:
		return d.FromString(casted)
	default:
		return fmt.Errorf("failed to convert driver value of type %T to Decimal", casted)
	}

	return nil
}

// MarshalJSON converts a Provider value to a JSON value
func (enum Provider) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Provider_name, ProviderMapping, true)), nil
}

// MarshalCSV converts a Provider value to CSV cell value
func (enum Provider) MarshalCSV() (string, error) {
	return utils.MarshalString(enum, Provider_name, ProviderMapping, false), nil
}

// MarshalYAML converts a Provider value to a YAML node value
func (enum Provider) MarshalYAML() (interface{}, error) {
	return utils.MarshalString(enum, Provider_name, ProviderMapping, false), nil
}

// MarshalDynamoDBAttributeValue converts a Provider value to a DynamoDB AttributeValue
func (enum Provider) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{Value: utils.MarshalString(enum, Provider_name, ProviderMapping, false)}, nil
}

// UnmarshalJSON attempts to convert a JSON value to a new Provider value
func (enum *Provider) UnmarshalJSON(raw []byte) error {
	return utils.UnmarshalValue(raw, Provider_value, ProviderAlternates, enum)
}

// UnmarshalCSV attempts to convert a CSV cell value to a new Provider value
func (enum *Provider) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Provider_value, ProviderAlternates, enum)
}

// UnmarshalYAML attempts to convert a YAML node to a new Provider value
func (enum *Provider) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.ScalarNode {
		return fmt.Errorf("YAML node had an invalid kind (expected scalar value)")
	} else {
		return utils.UnmarshalString(value.Value, Provider_value, ProviderAlternates, enum)
	}
}

// UnmarshalDynamoDBAttributeValue attempts to convert a DynamoDB AttributeVAlue to a Provider
// value. This function can handle []bytes, numerics, or strings. If the AttributeValue is NULL then
// the Provider value will not be modified.
func (enum *Provider) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Provider_value, ProviderAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Provider_value, ProviderAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Provider_value, ProviderAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Provider", value)
	}
}

// MarhsalJSON converts a Timestamp to JSON
func (timestamp *UnixTimestamp) MarshalJSON() ([]byte, error) {
	return []byte(timestamp.ToEpoch()), nil
}

// MarshalCSV converts a Timestamp to a CSV format
func (timestamp *UnixTimestamp) MarshalCSV() (string, error) {
	return timestamp.ToEpoch(), nil
}

// Marshaler converts a Timestamp to a DynamoDB attribute value
func (timestamp *UnixTimestamp) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: timestamp.ToEpoch(),
	}, nil
}

// Value converts a Timestamp to an SQL value
func (timestamp *UnixTimestamp) Value() (driver.Value, error) {
	return driver.Value(timestamp.ToEpoch()), nil
}

// UnmarshalJSON converts JSON data into a Timestamp
func (timestamp *UnixTimestamp) UnmarshalJSON(data []byte) error {

	// Check if the value is nil; if this is the case then return nil
	if data == nil {
		return nil
	}

	// Attempt to deserialize the value to a string to remove any escapes or
	// quotes that aren't needed; if this fails then return an error. If the
	// string isn't already quoted then we probably don't have any work to do
	// here so just trim the whitespace off and set it directly
	var asStr string
	if runes := []rune(string(data)); len(runes) >= 2 && runes[0] == '"' && runes[len(runes)-1] == '"' {
		if err := json.Unmarshal(data, &asStr); err != nil {
			return err
		}
	} else {
		asStr = string(data)
	}

	// Otherwise, convert the data from a string into a timestamp
	return timestamp.FromString(asStr)
}

// UnmarshalCSV converts a CSV column into a Timestamp
func (timestamp *UnixTimestamp) UnmarshalCSV(raw string) error {
	return timestamp.FromString(raw)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a timestamp
func (timestamp *UnixTimestamp) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return timestamp.FromString(string(casted.Value))
	case *types.AttributeValueMemberN:
		return timestamp.FromString(casted.Value)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return timestamp.FromString(casted.Value)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a UnixTimestamp", value)
	}
}

// Scan converts an SQL value into a Timestamp
func (timestamp *UnixTimestamp) Scan(value interface{}) error {

	// Check if the value is nil; if this is the case then return nil
	if value == nil {
		return nil
	}

	// Otherwise, convert the data from a string into a timestamp
	switch casted := value.(type) {
	case string:
		return timestamp.FromString(casted)
	case int64:
		timestamp.Seconds = casted / nanosPerSecond
		timestamp.Nanoseconds = int32(casted % nanosPerSecond)
		return nil
	default:
		return fmt.Errorf("Value of %v with a type of %T could not be converted to a UnixTimestamp", casted, casted)
	}
}

// MarhsalJSON converts a Duration to JSON
func (duration *UnixDuration) MarshalJSON() ([]byte, error) {
	return []byte(duration.ToEpoch()), nil
}

// MarshalCSV converts a Duration to a CSV format
func (duration *UnixDuration) MarshalCSV() (string, error) {
	return duration.ToEpoch(), nil
}

// Marshaler converts a Duration to a DynamoDB attribute value
func (duration *UnixDuration) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: duration.ToEpoch(),
	}, nil
}

// Value converts a Duration to an SQL value
func (duration *UnixDuration) Value() (driver.Value, error) {
	return driver.Value(duration.ToEpoch()), nil
}

// UnmarshalJSON converts JSON data into a Duration
func (duration *UnixDuration) UnmarshalJSON(data []byte) error {

	// Check if the value is nil; if this is the case then return nil
	if data == nil {
		return nil
	}

	// Attempt to deserialize the value to a string to remove any escapes or
	// quotes that aren't needed; if this fails then return an error. If the
	// string isn't already quoted then we probably don't have any work to do
	// here so just trim the whitespace off and set it directly
	var asStr string
	if runes := []rune(string(data)); len(runes) >= 2 && runes[0] == '"' && runes[len(runes)-1] == '"' {
		if err := json.Unmarshal(data, &asStr); err != nil {
			return err
		}
	} else {
		asStr = string(data)
	}

	// Otherwise, convert the data from a string into a duration
	return duration.FromString(asStr)
}

// UnmarshalCSV converts a CSV column into a Duration
func (duration *UnixDuration) UnmarshalCSV(raw string) error {
	return duration.FromString(raw)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Duration
func (duration *UnixDuration) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return duration.FromString(string(casted.Value))
	case *types.AttributeValueMemberN:
		return duration.FromString(casted.Value)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return duration.FromString(casted.Value)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a UnixDuration", value)
	}
}

// Scan converts an SQL value into a Duration
func (duration *UnixDuration) Scan(value interface{}) error {

	// Check if the value is nil; if this is the case then return nil
	if value == nil {
		return nil
	}

	// Otherwise, convert the data from a string into a duration
	return duration.FromString(value.(string))
}

// MarhsalJSON converts a Financial.Common.AssetClass to JSON
func (enum Financial_Common_AssetClass) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Common_AssetClass_name, AssetClassMapping, true)), nil
}

// MarshalCSV converts a Financial.Common.AssetClass to a CSV format
func (enum Financial_Common_AssetClass) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Common.AssetClass to a DynamoDB attribute value
func (enum Financial_Common_AssetClass) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Common_AssetClass_name, AssetClassMapping, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Common.AssetClass
func (enum *Financial_Common_AssetClass) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Common_AssetClass_value, AssetClassAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Common.AssetClass
func (enum *Financial_Common_AssetClass) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Common_AssetClass_value, AssetClassAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Common.AssetClass
func (enum *Financial_Common_AssetClass) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Common_AssetClass_value, AssetClassAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Common_AssetClass_value, AssetClassAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Common_AssetClass_value, AssetClassAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Common.AssetClass", value)
	}
}

// Scan converts an SQL value into a Financial.Common.AssetClass
func (enum *Financial_Common_AssetClass) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Common_AssetClass_value, AssetClassAlternates, enum)
}

// MarhsalJSON converts a Financial.Common.AssetType to JSON
func (enum Financial_Common_AssetType) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Common_AssetType_name, AssetTypeMapping, true)), nil
}

// MarshalCSV converts a Financial.Common.AssetType to a CSV format
func (enum Financial_Common_AssetType) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Common.AssetType to a DynamoDB attribute value
func (enum Financial_Common_AssetType) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Common_AssetType_name, AssetTypeMapping, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Common.AssetType
func (enum *Financial_Common_AssetType) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Common_AssetType_value, AssetTypeAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Common.AssetType
func (enum *Financial_Common_AssetType) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Common_AssetType_value, AssetTypeAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Common.AssetType
func (enum *Financial_Common_AssetType) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Common_AssetType_value, AssetTypeAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Common_AssetType_value, AssetTypeAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Common_AssetType_value, AssetTypeAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Common.AssetType", value)
	}
}

// Scan converts an SQL value into a Financial.Common.AssetType
func (enum *Financial_Common_AssetType) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Common_AssetType_value, AssetTypeAlternates, enum)
}

// MarhsalJSON converts a Financial.Common.Locale to JSON
func (enum Financial_Common_Locale) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Common_Locale_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Common.Locale to a CSV format
func (enum Financial_Common_Locale) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Common.Locale to a DynamoDB attribute value
func (enum Financial_Common_Locale) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Common_Locale_name, utils.Ignore, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Common.Locale
func (enum *Financial_Common_Locale) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Common_Locale_value, LocaleAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Common.Locale
func (enum *Financial_Common_Locale) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Common_Locale_value, LocaleAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Common.Locale
func (enum *Financial_Common_Locale) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Common_Locale_value, LocaleAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Common_Locale_value, LocaleAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Common_Locale_value, LocaleAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Common.Locale", value)
	}
}

// Scan converts an SQL value into a Financial.Common.Locale
func (enum *Financial_Common_Locale) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Common_Locale_value, LocaleAlternates, enum)
}

// MarhsalJSON converts a Financial.Common.Tape to JSON
func (enum Financial_Common_Tape) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Common_Tape_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Common.Tape to a CSV format
func (enum Financial_Common_Tape) MarshalCSV() (string, error) {
	return utils.MarshalString(enum, Financial_Common_Tape_name, utils.Ignore, false), nil
}

// Marshaler converts a Financial.Common.Tape to a DynamoDB attribute value
func (enum Financial_Common_Tape) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Common_Tape_name, utils.Ignore, false),
	}, nil
}

// Value converts a Financial.Common.Tape to an SQL value
func (enum Financial_Common_Tape) Value() (driver.Value, error) {
	return driver.Value(utils.MarshalString(enum, Financial_Common_Tape_name, utils.Ignore, false)), nil
}

// UnmarshalJSON converts JSON data into a Financial.Common.Tape
func (enum *Financial_Common_Tape) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Common_Tape_value, utils.None, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Common.Tape
func (enum *Financial_Common_Tape) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Common_Tape_value, utils.None, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Common.Tape
func (enum *Financial_Common_Tape) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Common_Tape_value, utils.None, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Common_Tape_value, utils.None, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Common_Tape_value, utils.None, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Common.Tape", value)
	}
}

// Scan converts an SQL value into a Financial.Common.Tape
func (enum *Financial_Common_Tape) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Common_Tape_value, utils.None, enum)
}

// MarhsalJSON converts a Financial.Dividends.Frequency to JSON
func (enum Financial_Dividends_Frequency) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Dividends_Frequency_name, DividendFrequencyMapping, true)), nil
}

// MarshalCSV converts a Financial.Dividends.Frequency to a CSV format
func (enum Financial_Dividends_Frequency) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Dividends.Frequency to a DynamoDB attribute value
func (enum Financial_Dividends_Frequency) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Dividends_Frequency_name, DividendFrequencyMapping, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Dividends.Frequency
func (enum *Financial_Dividends_Frequency) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Dividends_Frequency_value, DividendFrequencyAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Dividends.Frequency
func (enum *Financial_Dividends_Frequency) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Dividends_Frequency_value, DividendFrequencyAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Dividends.Frequency
func (enum *Financial_Dividends_Frequency) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Dividends_Frequency_value, DividendFrequencyAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Dividends_Frequency_value, DividendFrequencyAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Dividends_Frequency_value, DividendFrequencyAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Dividends.Frequency", value)
	}
}

// Scan converts an SQL value into a Financial.Dividends.Frequency
func (enum *Financial_Dividends_Frequency) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Dividends_Frequency_value, DividendFrequencyAlternates, enum)
}

// MarhsalJSON converts a Financial.Dividends.Type to JSON
func (enum Financial_Dividends_Type) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Dividends_Type_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Dividends.Type to a CSV format
func (enum Financial_Dividends_Type) MarshalCSV() (string, error) {
	return utils.MarshalString(enum, Financial_Dividends_Type_name, utils.Ignore, false), nil
}

// Marshaler converts a Financial.Dividends.Type to a DynamoDB attribute value
func (enum Financial_Dividends_Type) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Dividends_Type_name, utils.Ignore, false),
	}, nil
}

// Value converts a Financial.Dividends.Type to an SQL value
func (enum Financial_Dividends_Type) Value() (driver.Value, error) {
	return driver.Value(utils.MarshalString(enum, Financial_Dividends_Type_name, utils.Ignore, false)), nil
}

// UnmarshalJSON converts JSON data into a Financial.Dividends.Type
func (enum *Financial_Dividends_Type) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Dividends_Type_value, utils.None, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Dividends.Type
func (enum *Financial_Dividends_Type) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Dividends_Type_value, utils.None, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Dividends.Type
func (enum *Financial_Dividends_Type) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Dividends_Type_value, utils.None, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Dividends_Type_value, utils.None, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Dividends_Type_value, utils.None, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Dividends.Type", value)
	}
}

// Scan converts an SQL value into a Financial.Dividends.Type
func (enum *Financial_Dividends_Type) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Dividends_Type_value, utils.None, enum)
}

// MarhsalJSON converts a Financial.Exchanges.Type to JSON
func (enum Financial_Exchanges_Type) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Exchanges_Type_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Exchanges.Type to a CSV format
func (enum Financial_Exchanges_Type) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Exchanges.Type to a DynamoDB attribute value
func (enum Financial_Exchanges_Type) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Exchanges_Type_name, utils.Ignore, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Exchanges.Type
func (enum *Financial_Exchanges_Type) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Exchanges_Type_value, ExchangeTypeAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Exchanges.Type
func (enum *Financial_Exchanges_Type) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Exchanges_Type_value, ExchangeTypeAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Exchanges.Type
func (enum *Financial_Exchanges_Type) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Exchanges_Type_value, ExchangeTypeAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Exchanges_Type_value, ExchangeTypeAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Exchanges_Type_value, ExchangeTypeAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Exchanges.Type", value)
	}
}

// Scan converts an SQL value into a Financial.Exchanges.Type
func (enum *Financial_Exchanges_Type) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Exchanges_Type_value, ExchangeTypeAlternates, enum)
}

// MarhsalJSON converts a Financial.Options.ContractType to JSON
func (enum Financial_Options_ContractType) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Options_ContractType_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Options.ContractType to a CSV format
func (enum Financial_Options_ContractType) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Options.ContractType to a DynamoDB attribute value
func (enum Financial_Options_ContractType) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Options_ContractType_name, utils.Ignore, false),
	}, nil
}

// Value converts a Financial.Options.ContractType to an SQL value
func (enum Financial_Options_ContractType) Value() (driver.Value, error) {
	return driver.Value(utils.MarshalString(enum, Financial_Options_ContractType_name, utils.Ignore, false)), nil
}

// UnmarshalJSON converts JSON data into a Financial.Options.ContractType
func (enum *Financial_Options_ContractType) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Options_ContractType_value, OptionContractTypeAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Options.ContractType
func (enum *Financial_Options_ContractType) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Options_ContractType_value, OptionContractTypeAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Options.ContractType
func (enum *Financial_Options_ContractType) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Options_ContractType_value, OptionContractTypeAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Options_ContractType_value, OptionContractTypeAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Options_ContractType_value, OptionContractTypeAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Options.ContractType", value)
	}
}

// Scan converts an SQL value into a Financial.Options.ContractType
func (enum *Financial_Options_ContractType) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Options_ContractType_value, OptionContractTypeAlternates, enum)
}

// MarhsalJSON converts a Financial.Options.ExerciseStyle to JSON
func (enum Financial_Options_ExerciseStyle) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Options_ExerciseStyle_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Options.ExerciseStyle to a CSV format
func (enum Financial_Options_ExerciseStyle) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Options.ExerciseStyle to a DynamoDB attribute value
func (enum Financial_Options_ExerciseStyle) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Options_ExerciseStyle_name, utils.Ignore, false),
	}, nil
}

// Value converts a Financial.Options.ExerciseStyle to an SQL value
func (enum Financial_Options_ExerciseStyle) Value() (driver.Value, error) {
	return driver.Value(utils.MarshalString(enum, Financial_Options_ExerciseStyle_name, utils.Ignore, false)), nil
}

// UnmarshalJSON converts JSON data into a Financial.Options.ExerciseStyle
func (enum *Financial_Options_ExerciseStyle) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Options_ExerciseStyle_value, OptionExerciseStyleAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Options.ExerciseStyle
func (enum *Financial_Options_ExerciseStyle) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Options_ExerciseStyle_value, OptionExerciseStyleAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Options.ExerciseStyle
func (enum *Financial_Options_ExerciseStyle) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Options_ExerciseStyle_value, OptionExerciseStyleAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Options_ExerciseStyle_value, OptionExerciseStyleAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Options_ExerciseStyle_value, OptionExerciseStyleAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Options.ExerciseStyle", value)
	}
}

// Scan converts an SQL value into a Financial.Options.ExerciseStyle
func (enum *Financial_Options_ExerciseStyle) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Options_ExerciseStyle_value, OptionExerciseStyleAlternates, enum)
}

// MarhsalJSON converts a Financial.Options.UnderlyingType to JSON
func (enum Financial_Options_UnderlyingType) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Options_UnderlyingType_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Options.UnderlyingType to a CSV format
func (enum Financial_Options_UnderlyingType) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Options.UnderlyingType to a DynamoDB attribute value
func (enum Financial_Options_UnderlyingType) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Options_UnderlyingType_name, utils.Ignore, false),
	}, nil
}

// Value converts a Financial.Options.UnderlyingType to an SQL value
func (enum Financial_Options_UnderlyingType) Value() (driver.Value, error) {
	return driver.Value(utils.MarshalString(enum, Financial_Options_UnderlyingType_name, utils.Ignore, false)), nil
}

// UnmarshalJSON converts JSON data into a Financial.Options.UnderlyingType
func (enum *Financial_Options_UnderlyingType) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Options_UnderlyingType_value, OptionUnderlyingTypeAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Options.UnderlyingType
func (enum *Financial_Options_UnderlyingType) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Options_UnderlyingType_value, OptionUnderlyingTypeAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Options.UnderlyingType
func (enum *Financial_Options_UnderlyingType) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Options_UnderlyingType_value, OptionUnderlyingTypeAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Options_UnderlyingType_value, OptionUnderlyingTypeAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Options_UnderlyingType_value, OptionUnderlyingTypeAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Options.UnderlyingType", value)
	}
}

// Scan converts an SQL value into a Financial.Options.UnderlyingType
func (enum *Financial_Options_UnderlyingType) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Options_UnderlyingType_value, OptionUnderlyingTypeAlternates, enum)
}

// MarhsalJSON converts a Financial.Quotes.Condition to JSON
func (enum Financial_Quotes_Condition) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Quotes_Condition_name, QuoteConditionMapping, true)), nil
}

// MarshalCSV converts a Financial.Quotes.Condition to a CSV format
func (enum Financial_Quotes_Condition) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Quotes.Condition to a DynamoDB attribute value
func (enum Financial_Quotes_Condition) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Quotes_Condition_name, QuoteConditionMapping, false),
	}, nil
}

// Value converts a Financial.Quotes.Condition to an SQL value
func (enum Financial_Quotes_Condition) Value() (driver.Value, error) {

	// If we have an invalid value then return the actual value for invalid
	if enum == Financial_Quotes_Invalid {
		return driver.Value(-1), nil
	}

	// Otherwise, let the driver use the integer value of the enum
	return driver.Value(int(enum)), nil
}

// UnmarshalJSON converts JSON data into a Financial.Quotes.Condition
func (enum *Financial_Quotes_Condition) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Quotes_Condition_value, QuoteConditionAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Quotes.Condition
func (enum *Financial_Quotes_Condition) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Quotes_Condition_value, QuoteConditionAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Quotes.Condition
func (enum *Financial_Quotes_Condition) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Quotes_Condition_value, QuoteConditionAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Quotes_Condition_value, QuoteConditionAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Quotes_Condition_value, QuoteConditionAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Quotes.Condition", value)
	}
}

// Scan converts an SQL value into a Financial.Quotes.Condition
func (enum *Financial_Quotes_Condition) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Quotes_Condition_value, QuoteConditionAlternates, enum)
}

// MarhsalJSON converts a Financial.Quotes.Indicator to JSON
func (enum Financial_Quotes_Indicator) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Quotes_Indicator_name, QuoteIndicatorMapping, true)), nil
}

// MarshalCSV converts a Financial.Quotes.Indicator to a CSV format
func (enum Financial_Quotes_Indicator) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Quotes.Indicator to a DynamoDB attribute value
func (enum Financial_Quotes_Indicator) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Quotes_Indicator_name, QuoteIndicatorMapping, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Quotes.Indicator
func (enum *Financial_Quotes_Indicator) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Quotes_Indicator_value, QuoteIndicatorAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Quotes.Indicator
func (enum *Financial_Quotes_Indicator) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Quotes_Indicator_value, QuoteIndicatorAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Quotes.Indicator
func (enum *Financial_Quotes_Indicator) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Quotes_Indicator_value, QuoteIndicatorAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Quotes_Indicator_value, QuoteIndicatorAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Quotes_Indicator_value, QuoteIndicatorAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Quotes.Indicator", value)
	}
}

// Scan converts an SQL value into a Financial.Quotes.Indicator
func (enum *Financial_Quotes_Indicator) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Quotes_Indicator_value, QuoteIndicatorAlternates, enum)
}

// MarhsalJSON converts a Financial.Trades.Condition to JSON
func (enum Financial_Trades_Condition) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Trades_Condition_name, TradeConditionMapping, true)), nil
}

// MarshalCSV converts a Financial.Trades.Condition to a CSV format
func (enum Financial_Trades_Condition) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Trades.Condition to a DynamoDB attribute value
func (enum Financial_Trades_Condition) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Trades_Condition_name, TradeConditionMapping, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Trades.Condition
func (enum *Financial_Trades_Condition) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Trades_Condition_value, TradeConditionAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Trades.Condition
func (enum *Financial_Trades_Condition) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Trades_Condition_value, TradeConditionAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Trades.Condition
func (enum *Financial_Trades_Condition) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Trades_Condition_value, TradeConditionAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Trades_Condition_value, TradeConditionAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Trades_Condition_value, TradeConditionAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Trades.Condition", value)
	}
}

// Scan converts an SQL value into a Financial.Trades.Condition
func (enum *Financial_Trades_Condition) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Trades_Condition_value, TradeConditionAlternates, enum)
}

// MarhsalJSON converts a Financial.Trades.CorrectionCode to JSON
func (enum Financial_Trades_CorrectionCode) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Trades_CorrectionCode_name, TradeCorrectionMapping, true)), nil
}

// MarshalCSV converts a Financial.Trades.CorrectionCode to a CSV format
func (enum Financial_Trades_CorrectionCode) MarshalCSV() (string, error) {
	return fmt.Sprintf("%02d", enum), nil
}

// Marshaler converts a Financial.Trades.CorrectionCode to a DynamoDB attribute value
func (enum Financial_Trades_CorrectionCode) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Trades_CorrectionCode_name, TradeCorrectionMapping, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Trades.CorrectionCode
func (enum *Financial_Trades_CorrectionCode) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Trades_CorrectionCode_value, TradeCorrectionAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Trades.CorrectionCode
func (enum *Financial_Trades_CorrectionCode) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Trades_CorrectionCode_value, TradeCorrectionAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Trades.CorrectionCode
func (enum *Financial_Trades_CorrectionCode) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Trades_CorrectionCode_value, TradeCorrectionAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Trades_CorrectionCode_value, TradeCorrectionAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Trades_CorrectionCode_value, TradeCorrectionAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Trades.CorrectionCode", value)
	}
}

// Scan converts an SQL value into a Financial.Trades.CorrectionCode
func (enum *Financial_Trades_CorrectionCode) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Trades_CorrectionCode_value, TradeCorrectionAlternates, enum)
}
