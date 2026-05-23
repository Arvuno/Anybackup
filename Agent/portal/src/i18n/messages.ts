export type Locale = "en" | "zh-CN"

export type MessageKey =
  | "common.appName"
  | "common.appTagline"
  | "common.loading"
  | "login.username"
  | "login.password"
  | "login.usernamePlaceholder"
  | "login.passwordPlaceholder"
  | "login.showPassword"
  | "login.hidePassword"
  | "login.submit"
  | "login.submitting"
  | "login.serviceUnavailable"
  | "validation.requiredUsername"
  | "validation.requiredPassword"
  | "validation.passwordMinLength"
  | "validation.passwordSameAsUsername"
  | "validation.passwordNeedsLetterAndDigit"
  | "validation.passwordConfirmationMismatch"
  | "sidebar.demoUser"
  | "sidebar.backupAdmin"
  | "sidebar.loginAccount"
  | "sidebar.logout"
  | "sidebar.newConversation"
  | "sidebar.searchConversation"
  | "sidebar.refreshingConversations"
  | "sidebar.noMatchedConversation"
  | "sidebar.noConversationYet"
  | "sidebar.historyConversation"
  | "sidebar.expandSidebar"
  | "sidebar.collapseSidebar"
  | "sidebar.settings"
  | "settings.title"
  | "settings.description"
  | "settings.userManagementTitle"
  | "settings.userManagementDescription"
  | "guide.progressLabel"
  | "guide.startLearn"
  | "guide.howItWorksTitle"
  | "guide.howItWorksDescription"
  | "guide.viewDemo"
  | "guide.userDescribeTitle"
  | "guide.userDescribeDescription"
  | "guide.agentGenerateTitle"
  | "guide.agentGenerateDescription"
  | "guide.pendingReview"
  | "guide.nextReview"
  | "guide.agentGeneratingPlan"
  | "guide.reviewTitle"
  | "guide.reviewDescription"
  | "guide.reviewPassed"
  | "guide.planApproved"
  | "guide.agentExecuteBackup"
  | "guide.finishGuide"
  | "guide.readyTitle"
  | "guide.readyDescription"
  | "guide.enterWorkspace"
  | "guide.dontShowAgain"
  | "guide.skipGuide"
  | "guide.workflowStep1Label"
  | "guide.workflowStep1Sub"
  | "guide.workflowStep2Label"
  | "guide.workflowStep2Sub"
  | "guide.workflowStep3Label"
  | "guide.workflowStep3Sub"
  | "guide.workflowStep4Label"
  | "guide.workflowStep4Sub"
  | "guide.product1Desc"
  | "guide.product2Desc"
  | "guide.product3Desc"
  | "guide.feature1"
  | "guide.feature2"
  | "guide.feature3"
  | "guide.feature4"
  | "guide.feature5"
  | "guide.userText"
  | "guide.agentText"
  | "guide.planTitle"
  | "guide.planEnvironment"
  | "guide.planMetricBackupType"
  | "guide.planMetricStrategy"
  | "guide.planMetricRpoRto"
  | "guide.planMetricDataReduction"
  | "guide.planMetricBackupTypeValue"
  | "guide.planMetricStrategyValue"
  | "guide.planMetricRpoRtoValue"
  | "guide.planMetricDataReductionValue"
  | "guide.planRecommendation"
  | "guide.reviewPlanSummary"
  | "validation.requiredDisplayName"
  | "validation.confirmPasswordRequired"
  | "errors.genericUnavailable"
  | "errors.networkUnavailable"
  | "conversation.placeholder.richContent"
  | "conversation.placeholder.clarification"
  | "conversation.placeholder.status"
  | "conversation.error.listLoadFailed"
  | "conversation.error.detailLoadFailed"
  | "conversation.error.messagesLoadFailed"
  | "conversation.error.eventsLoadFailed"
  | "conversation.error.createFailed"
  | "conversation.error.sendFailed"
  | "conversation.error.archiveFailed"
  | "conversation.error.restoreFailed"
  | "conversation.error.copyConfigFailed"
  | "conversation.pendingTurnTimeout"
  | "conversation.serviceUnavailable"
  | "conversation.newConversationTitle"
  | "auth.authServiceUnavailable"
  | "auth.userDisabled"
  | "auth.invalidCredentials"
  | "auth.sessionValidationFailed"
  | "auth.sessionExpired"
  | "auth.logoutFailed"
  | "users.unnamed"
  | "users.listLoadFailed"
  | "users.detailLoadFailed"
  | "users.roleLoadFailed"
  | "users.adminRoleNotFound"
  | "users.assignRoleFailed"
  | "users.createReadFailed"
  | "users.createFailed"
  | "users.assignAfterCreateFailed"
  | "users.usernameExists"
  | "users.saveFailed"
  | "users.statusSaveFailed"
  | "users.cannotDisableSelf"
  | "users.resetPasswordFailed"
  | "users.feedback.created"
  | "users.feedback.saved"
  | "users.feedback.enabled"
  | "users.feedback.disabled"
  | "users.feedback.passwordReset"
  | "chat.closeHint"
  | "chat.preparingWorkspace"
  | "chat.restoringConversation"
  | "chat.refNotInMessagePrefix"
  | "chat.refreshingList"
  | "chat.conversationNotLoaded"
  | "chat.composer.placeholder"
  | "chat.composer.hint"
  | "chat.composer.send"
  | "chat.composer.inputLabel"
  | "chat.empty.title"
  | "chat.empty.description"
  | "chat.panelHeader.title"
  | "chat.panelHeader.subtitle"
  | "chat.waiting.thinking"
  | "chat.waiting.duration"
  | "chat.thought.label"
  | "chat.thought.running"
  | "chat.thought.done"
  | "chat.messageTime.invalid"
  | "candidate.recommended"
  | "candidate.confirmed"
  | "workspace.breadcrumb.home"
  | "workspace.breadcrumb.settings"
  | "workspace.breadcrumb.users"
  | "workspace.navTitle"
  | "modal.close"
  | "language.zhShort"
  | "usersPage.title"
  | "usersPage.subtitle"
  | "usersPage.enabledCount"
  | "usersPage.createUser"
  | "usersPage.listTitle"
  | "usersPage.listHint"
  | "usersPage.userCount"
  | "usersPage.colUsername"
  | "usersPage.colDisplayName"
  | "usersPage.colRole"
  | "usersPage.colStatus"
  | "usersPage.colActions"
  | "usersPage.loading"
  | "usersPage.empty"
  | "usersPage.currentAccount"
  | "usersPage.edit"
  | "usersPage.disable"
  | "usersPage.enable"
  | "usersPage.resetPassword"
  | "userForm.titleCreate"
  | "userForm.titleEdit"
  | "userForm.descCreate"
  | "userForm.descEdit"
  | "userForm.cancel"
  | "userForm.create"
  | "userForm.save"
  | "userForm.labelUsername"
  | "userForm.labelDisplayName"
  | "userForm.labelPassword"
  | "userForm.labelConfirmPassword"
  | "userForm.status"
  | "userForm.statusHint"
  | "userForm.statusEnabled"
  | "userForm.statusDisabled"
  | "userForm.enableToDisableHint"
  | "userForm.disableToEnableHint"
  | "userForm.formError"
  | "resetPassword.title"
  | "resetPassword.newPassword"
  | "resetPassword.confirmNewPassword"
  | "resetPassword.description"
  | "resetPassword.cancel"
  | "resetPassword.submit"
  | "resetPassword.targetLabel"
  | "disableUser.title"
  | "disableUser.description"
  | "disableUser.cancel"
  | "disableUser.confirm"
  | "disableUser.warningTitle"
  | "disableUser.warningBody"
  | "disableUser.preview"
  | "disableUser.selectUser"
  | "demo.pageTitle"
  | "demo.pageSubtitle"
  | "demo.visitPath"
  | "demo.sectionLayoutTreeDesc"
  | "demo.protocolActionFeedback"
  | "demo.sectionLayoutTitle"
  | "demo.sectionNodeGalleryDesc"
  | "demo.sectionCandidateTitle"
  | "demo.sectionCandidateDesc"
  | "demo.lastInteraction"
  | "demo.sectionTextTitle"
  | "demo.sectionTextDesc"
  | "demo.sectionWaitingTitle"
  | "demo.sectionWaitingDesc"
  | "demo.listTitleRestore"
  | "demo.listSummaryCandidate"
  | "demo.listTitleBackupQuery"
  | "demo.listSummaryPlain"
  | "demo.listTitleGenerating"
  | "demo.listSummaryWaiting"
  | "demo.selection.none"
  | "demo.selection.confirm"
  | "demo.selection.reject"
  | "demo.selection.revise"
  | "demo.selection.reviseEmpty"
  | "demo.layout.none"
  | "demo.layout.last"
  | "demo.text.user.backupStatus"
  | "demo.text.assistant.backupStatus"
  | "demo.waiting.user.restorePlan"
  | "demo.candidate.user.message"
  | "demo.candidate.assistant.lead"
  | "demo.candidate.cardTitle"
  | "demo.candidate.cardSummary"
  | "demo.candidate.actionConfirm"
  | "demo.candidate.actionReject"
  | "demo.candidate.actionRevise"
  | "demo.candidate.reviseInputLabel"
  | "demo.candidate.revisePlaceholder"
  | "demo.candidate.reviseSubmit"
  | "demo.optionA.title"
  | "demo.optionA.summary"
  | "demo.optionA.extraTitle"
  | "demo.optionA.extraBody"
  | "demo.optionB.title"
  | "demo.optionB.summary"
  | "demo.optionB.extraTitle"
  | "demo.optionB.extraBody"
  | "demo.optionC.title"
  | "demo.optionC.summary"
  | "demo.optionC.extraTitle"
  | "demo.optionC.extraBody"
  | "demo.field.restoreScope"
  | "demo.field.rpoRto"
  | "demo.field.target"
  | "demo.field.coverage"
  | "demo.value.dbAsset"
  | "demo.value.rpoRtoA"
  | "demo.value.targetMysql"
  | "demo.value.coverageTable"
  | "demo.value.rpoRtoB"
  | "demo.value.targetProd"
  | "demo.value.coverageFullDb"
  | "demo.value.rpoRtoC"
  | "demo.value.targetIso"
  | "demo.value.coverageValidate"
  | "demo.layout.heading"
  | "demo.layout.intro"
  | "demo.layout.badgeRecommended"
  | "demo.layout.badgeMediumRisk"
  | "demo.layout.paraA"
  | "demo.metric.rpo"
  | "demo.metric.rto"
  | "demo.metric.success"
  | "demo.metric.valueRpo"
  | "demo.metric.valueRto"
  | "demo.metric.valueHigh"
  | "demo.layout.kvFullDb"
  | "demo.layout.kvResourceHigh"
  | "demo.layout.kvComplexityMid"
  | "demo.layout.resourceUsageLabel"
  | "demo.layout.execComplexityLabel"
  | "demo.layout.badgeHighResource"
  | "demo.layout.paraB"
  | "demo.layout.titleB"
  | "demo.layout.badgeConsistency"
  | "demo.layout.paraC"
  | "demo.layout.titleC"
  | "demo.layout.riskTitle"
  | "demo.layout.riskBody"
  | "demo.layout.calloutBizTitle"
  | "demo.layout.calloutBizText"
  | "demo.action.confirmA"
  | "demo.action.viewDetail"
  | "demo.action.confirmB"
  | "demo.action.confirmC"

type MessageBundle = Record<MessageKey, string>

export const LANGUAGE_STORAGE_KEY = "agent.portal.locale"

const DEFAULT_LOCALE: Locale = "en"

export function getStoredLocale(): Locale {
  if (typeof window === "undefined") return DEFAULT_LOCALE
  const value = window.localStorage.getItem(LANGUAGE_STORAGE_KEY)
  if (value === "en" || value === "zh-CN") {
    return value
  }
  return DEFAULT_LOCALE
}

export const messages: Record<Locale, MessageBundle> = {
  en: {
    "common.appName": "Anybackup",
    "common.appTagline": "Always Resilient",
    "common.loading": "Loading...",
    "login.username": "Username",
    "login.password": "Password",
    "login.usernamePlaceholder": "Enter username",
    "login.passwordPlaceholder": "Enter password",
    "login.showPassword": "Show password",
    "login.hidePassword": "Hide password",
    "login.submit": "Sign in",
    "login.submitting": "Signing in",
    "login.serviceUnavailable": "Login service is temporarily unavailable. Please try again later.",
    "validation.requiredUsername": "Please enter username",
    "validation.requiredPassword": "Please enter password",
    "validation.passwordMinLength": "Password must be at least 8 characters",
    "validation.passwordSameAsUsername": "Password cannot be exactly the same as username",
    "validation.passwordNeedsLetterAndDigit": "Password must contain at least letters and numbers",
    "validation.passwordConfirmationMismatch": "The two passwords do not match",
    "sidebar.demoUser": "Demo User",
    "sidebar.backupAdmin": "Backup Administrator",
    "sidebar.loginAccount": "Account",
    "sidebar.logout": "Sign out",
    "sidebar.newConversation": "New conversation",
    "sidebar.searchConversation": "Search conversations",
    "sidebar.refreshingConversations": "Refreshing conversation list...",
    "sidebar.noMatchedConversation": "No matching conversations. Try another keyword.",
    "sidebar.noConversationYet": "No formal conversation yet. Start with a new one.",
    "sidebar.historyConversation": "Conversation history",
    "sidebar.expandSidebar": "Expand sidebar",
    "sidebar.collapseSidebar": "Collapse sidebar",
    "sidebar.settings": "Settings",
    "settings.title": "Settings",
    "settings.description": "Manage platform capabilities and user accounts.",
    "settings.userManagementTitle": "User Management",
    "settings.userManagementDescription":
      "View users, create users, enable or disable accounts, and reset passwords as an administrator.",
    "guide.progressLabel": "Guide progress",
    "guide.startLearn": "Get started",
    "guide.howItWorksTitle": "Workflow",
    "guide.howItWorksDescription": "Natural language driven, AI-generated plans, and human approval",
    "guide.viewDemo": "View demo",
    "guide.userDescribeTitle": "User describes requirements",
    "guide.userDescribeDescription": "Describe your backup needs in natural language, like talking to an expert",
    "guide.agentGenerateTitle": "Agent generates a plan",
    "guide.agentGenerateDescription": "AI analyzes environment and generates an optimal backup plan",
    "guide.pendingReview": "Pending review",
    "guide.nextReview": "Next: review plan",
    "guide.agentGeneratingPlan": "Agent is analyzing the environment and generating a plan...",
    "guide.reviewTitle": "AI creates, humans decide",
    "guide.reviewDescription": "All critical operations require your final confirmation",
    "guide.reviewPassed": "Approve plan",
    "guide.planApproved": "Plan approved",
    "guide.agentExecuteBackup": "Agent will now start executing the backup tasks automatically",
    "guide.finishGuide": "Finish guide",
    "guide.readyTitle": "All set!",
    "guide.readyDescription":
      "The three-pane workspace lets you observe assets, chat with Agent, and review plans at once.",
    "guide.enterWorkspace": "Enter workspace",
    "guide.dontShowAgain": "Do not show this next time",
    "guide.skipGuide": "Skip guide and enter now",
    "guide.workflowStep1Label": "Describe requirements naturally",
    "guide.workflowStep1Sub": "AI understands your resilience intent",
    "guide.workflowStep2Label": "Agent analyzes environment",
    "guide.workflowStep2Sub": "Automatic scan and optimal strategy recommendation",
    "guide.workflowStep3Label": "Review plan",
    "guide.workflowStep3Sub": "AI generates, humans decide",
    "guide.workflowStep4Label": "Execute and verify automatically",
    "guide.workflowStep4Sub": "Continuous 24x7 protection and inspection",
    "guide.product1Desc": "AI-native autonomous backup and recovery agent",
    "guide.product2Desc": "Backup and recovery platform for autonomous agent operations",
    "guide.product3Desc": "Asset backup client",
    "guide.feature1": "Autonomous awareness",
    "guide.feature2": "Zero barrier",
    "guide.feature3": "24x7 online",
    "guide.feature4": "Transparent and trusted",
    "guide.feature5": "Self-evolving",
    "guide.userText": "The ERP Oracle database needs backup protection",
    "guide.agentText":
      "Understood. I detected ERP Oracle 19c RAC in your environment.\n\nWhat are your RPO/RTO targets and available backup windows?",
    "guide.planTitle": "ERP Oracle RAC Backup Plan",
    "guide.planEnvironment": "Oracle 19c RAC",
    "guide.planMetricBackupType": "Backup Type",
    "guide.planMetricStrategy": "Strategy",
    "guide.planMetricRpoRto": "RPO / RTO",
    "guide.planMetricDataReduction": "Data Reduction",
    "guide.planMetricBackupTypeValue": "RMAN Full+Incr+Log",
    "guide.planMetricStrategyValue": "Daily full at 02:00",
    "guide.planMetricRpoRtoValue": "<5min / <30min",
    "guide.planMetricDataReductionValue": "Source-side dedup",
    "guide.planRecommendation": "This plan beats peer systems by 20% in RTO and is recommended.",
    "guide.reviewPlanSummary": "RMAN Full+Incr+Log · RPO<5min · Source-side dedup",
    "validation.requiredDisplayName": "Please enter display name",
    "validation.confirmPasswordRequired": "Please confirm password",
    "errors.genericUnavailable": "Service is temporarily unavailable. Please try again later.",
    "errors.networkUnavailable": "Network is unavailable. Check your connection and try again.",
    "conversation.placeholder.richContent": "[Structured content]",
    "conversation.placeholder.clarification": "[Pending confirmation]",
    "conversation.placeholder.status": "[Status update]",
    "conversation.error.listLoadFailed": "Failed to load conversations. Please try again later.",
    "conversation.error.detailLoadFailed": "Failed to load conversation details. Please try again later.",
    "conversation.error.messagesLoadFailed": "Failed to load messages. Please try again later.",
    "conversation.error.eventsLoadFailed": "Failed to load status events. Please try again later.",
    "conversation.error.createFailed": "Failed to create conversation. Please try again later.",
    "conversation.error.sendFailed": "Failed to send message. Please try again later.",
    "conversation.error.archiveFailed": "Failed to archive conversation. Please try again later.",
    "conversation.error.restoreFailed": "Failed to restore conversation. Please try again later.",
    "conversation.error.copyConfigFailed": "Failed to copy conversation settings. Please try again later.",
    "conversation.pendingTurnTimeout": "The request timed out. Please try again later.",
    "conversation.serviceUnavailable": "Conversation service is temporarily unavailable. Please try again later.",
    "conversation.newConversationTitle": "New conversation",
    "auth.authServiceUnavailable": "Authentication service is temporarily unavailable. Please try again later.",
    "auth.userDisabled": "Account is disabled. Contact an administrator.",
    "auth.invalidCredentials": "Invalid username or password.",
    "auth.sessionValidationFailed": "Session validation failed. Please sign in again.",
    "auth.sessionExpired": "Your session has expired. Please sign in again.",
    "auth.logoutFailed": "Sign-out failed. Please try again later.",
    "users.unnamed": "Unnamed user",
    "users.listLoadFailed": "Failed to load user list. Please try again later.",
    "users.detailLoadFailed": "Failed to load user profile. Please try again later.",
    "users.roleLoadFailed": "Failed to load built-in administrator role details. Please try again later.",
    "users.adminRoleNotFound":
      "No built-in administrator role available for assignment (tried: {{roles}}). Check role configuration in this environment.",
    "users.assignRoleFailed": "Failed to assign built-in administrator role. Please try again later.",
    "users.createReadFailed": "User was created but details could not be read.",
    "users.createFailed": "Failed to create user. Please try again later.",
    "users.assignAfterCreateFailed":
      "User was created but assigning the built-in administrator role failed. Please try again later.",
    "users.usernameExists": "Username already exists.",
    "users.saveFailed": "Failed to save user. Please try again later.",
    "users.statusSaveFailed": "Failed to save user status. Please try again later.",
    "users.cannotDisableSelf": "You cannot disable the currently signed-in user.",
    "users.resetPasswordFailed": "Failed to reset password. Please try again later.",
    "users.feedback.created": "User created.",
    "users.feedback.saved": "User saved.",
    "users.feedback.enabled": "User enabled.",
    "users.feedback.disabled": "User disabled.",
    "users.feedback.passwordReset": "Password reset. Ask the user to sign in with the new password.",
    "chat.closeHint": "Dismiss notice",
    "chat.preparingWorkspace": "Preparing conversation workspace...",
    "chat.restoringConversation": "Restoring conversation...",
    "chat.refNotInMessagePrefix": "Referenced content was not returned with the message and cannot be opened yet:",
    "chat.refreshingList": "Refreshing conversation list...",
    "chat.conversationNotLoaded": "Conversation is missing or not loaded. Select again.",
    "chat.composer.placeholder": "Tell me what you want to back up, restore, or check...",
    "chat.composer.hint": "Continue in this conversation; candidate confirmations will flow as structured input.",
    "chat.composer.send": "Send message",
    "chat.composer.inputLabel": "Message input",
    "chat.empty.title": "Start a new conversation here",
    "chat.empty.description": "Ask follow-ups in this conversation; candidate confirmations will proceed in the same thread.",
    "chat.panelHeader.title": "Conversation",
    "chat.panelHeader.subtitle": "Conversation flow will be designed in a follow-up plan.",
    "chat.waiting.thinking": "Thinking",
    "chat.waiting.duration": "Thought for {minutes}m{seconds}s",
    "chat.thought.label": "Thought",
    "chat.thought.running": "Thinking",
    "chat.thought.done": "Done",
    "chat.messageTime.invalid": "Just now",
    "candidate.recommended": "Recommended",
    "candidate.confirmed": "Confirmed",
    "workspace.breadcrumb.home": "Workspace",
    "workspace.breadcrumb.settings": "Settings",
    "workspace.breadcrumb.users": "User management",
    "workspace.navTitle": "Home",
    "modal.close": "Close",
    "language.zhShort": "中文",
    "usersPage.title": "User management",
    "usersPage.subtitle":
      "Manage sign-in accounts, initial passwords, and enable/disable state. New users receive the administrator role automatically.",
    "usersPage.enabledCount": "Enabled {enabled} / {total}",
    "usersPage.createUser": "New user",
    "usersPage.listTitle": "Users",
    "usersPage.listHint": "This page only shows sign-in account administration.",
    "usersPage.userCount": "{count} users",
    "usersPage.colUsername": "Username",
    "usersPage.colDisplayName": "Display name",
    "usersPage.colRole": "Role",
    "usersPage.colStatus": "Status",
    "usersPage.colActions": "Actions",
    "usersPage.loading": "Loading users...",
    "usersPage.empty": "No users yet",
    "usersPage.currentAccount": "Current account",
    "usersPage.edit": "Edit",
    "usersPage.disable": "Disable",
    "usersPage.enable": "Enable",
    "usersPage.resetPassword": "Reset password",
    "userForm.titleCreate": "New user",
    "userForm.titleEdit": "Edit user",
    "userForm.descCreate":
      "Administrator sets the initial password; the administrator role is assigned automatically after creation.",
    "userForm.descEdit": "Username cannot be changed after the user is created.",
    "userForm.cancel": "Cancel",
    "userForm.create": "Create user",
    "userForm.save": "Save",
    "userForm.labelUsername": "Username",
    "userForm.labelDisplayName": "Display name",
    "userForm.labelPassword": "Password",
    "userForm.labelConfirmPassword": "Confirm password",
    "userForm.status": "Status",
    "userForm.statusHint": "Controls whether this user may sign in.",
    "userForm.statusEnabled": "Enabled",
    "userForm.statusDisabled": "Disabled",
    "userForm.enableToDisableHint": "Currently enabled — click to disable",
    "userForm.disableToEnableHint": "Currently disabled — click to enable",
    "userForm.formError": "Service is temporarily unavailable. Please try again later.",
    "resetPassword.title": "Reset password",
    "resetPassword.newPassword": "New password",
    "resetPassword.confirmNewPassword": "Confirm new password",
    "resetPassword.description": "Set a new password for {{username}}. Passwords are not emailed automatically.",
    "resetPassword.cancel": "Cancel",
    "resetPassword.submit": "Reset password",
    "resetPassword.targetLabel": "Reset target",
    "disableUser.title": "Disable user",
    "disableUser.description":
      "After disabling, the user cannot sign in. History and profile data are retained.",
    "disableUser.cancel": "Cancel",
    "disableUser.confirm": "Confirm disable",
    "disableUser.warningTitle": "Disabling removes access immediately.",
    "disableUser.warningBody": "Re-enable the account later from user management if needed.",
    "disableUser.preview": "You are about to disable user {{username}}",
    "disableUser.selectUser": "Select a user to disable.",
    "demo.pageTitle": "Conversation state demo",
    "demo.pageSubtitle":
          "This page uses local fixtures only and reuses chat components to preview plain messages, waiting states, and candidate cards.",
    "demo.visitPath": "Path:",
    "demo.sectionLayoutTreeDesc":
      "Render the candidate comparison block with the layout-tree registry and action dispatcher per protocol.",
    "demo.protocolActionFeedback": "Protocol action feedback:",
    "demo.sectionLayoutTitle": "AG-UI Layout Tree Demo",
    "demo.sectionNodeGalleryDesc":
      "Show remaining whitelisted node types: tabs, markdown, table, chart, attachment, and actions.",
    "demo.sectionCandidateTitle": "Pattern 1 · Candidate confirmation card",
    "demo.sectionCandidateDesc":
      "Primary structured rich content scenario. Actions are handled locally to mimic the full loop.",
    "demo.lastInteraction": "Latest interaction:",
    "demo.sectionTextTitle": "Pattern 2 · Plain message",
    "demo.sectionTextDesc": "Typical Q&A text with header and timestamps.",
    "demo.sectionWaitingTitle": "Pattern 3 · Waiting for reply",
    "demo.sectionWaitingDesc":
      "User message is sent while the assistant is still processing; a waiting placeholder appears at the bottom.",
    "demo.listTitleRestore": "Order database restore",
    "demo.listSummaryCandidate": "Candidate card demo",
    "demo.listTitleBackupQuery": "Backup status query",
    "demo.listSummaryPlain": "Plain text message demo",
    "demo.listTitleGenerating": "Generating restore plan",
    "demo.listSummaryWaiting": "After send, waiting for the model",
    "demo.selection.none": "No interaction yet — use buttons on the candidate card to see feedback.",
    "demo.selection.confirm": "Simulated confirm:",
    "demo.selection.reject": "Simulated reject:",
    "demo.selection.revise": "Simulated additional constraints:",
    "demo.selection.reviseEmpty": "not provided",
    "demo.layout.none":
      "No action yet — use Confirm on a card to see how payloads surface in the UI.",
    "demo.layout.last": "Last action:",
    "demo.text.user.backupStatus": "Can you check the order database backup from early this morning?",
    "demo.text.assistant.backupStatus":
      "The 02:00 order DB backup finished successfully in 7m 12s. No issues detected; you can review restore points or export details.",
    "demo.waiting.user.restorePlan": "Please generate a plan to restore the order database to yesterday afternoon.",
    "demo.candidate.user.message": "Restore to 2026-04-19 afternoon — show candidate options first.",
    "demo.candidate.assistant.lead": "Found 3 usable restore points.",
    "demo.candidate.cardTitle": "3 candidate options for you",
    "demo.candidate.cardSummary":
      "Filtered 3 candidates by restore scope, business impact, and execution cost.",
    "demo.candidate.actionConfirm": "Confirm and submit",
    "demo.candidate.actionReject": "Discard option",
    "demo.candidate.actionRevise": "Add constraints",
    "demo.candidate.reviseInputLabel": "Additional constraints",
    "demo.candidate.revisePlaceholder": "e.g. Generate the plan only; do not execute yet.",
    "demo.candidate.reviseSubmit": "Submit constraints",
    "demo.optionA.title": "Option A · Cross-host DB restore + table export/import",
    "demo.optionA.summary": "Recommended — tight scope and shorter business window.",
    "demo.optionA.extraTitle": "Business impact",
    "demo.optionA.extraBody":
      "Restore scope is database-level but only the target table is touched, reducing risk to other tables.",
    "demo.optionB.title": "Option B · Direct database restore to production",
    "demo.optionB.summary": "Shortest path but impacts the live database.",
    "demo.optionB.extraTitle": "Business impact",
    "demo.optionB.extraBody": "Direct restore may overwrite other healthy tables — higher risk.",
    "demo.optionC.title": "Option C · Cross-host validation restore in a delayed window",
    "demo.optionC.summary": "Best when integrity must be validated first.",
    "demo.optionC.extraTitle": "Business impact",
    "demo.optionC.extraBody": "Takes longer but allows consistency checks before production restore.",
    "demo.field.restoreScope": "Restore scope",
    "demo.field.rpoRto": "RPO / RTO",
    "demo.field.target": "Target",
    "demo.field.coverage": "Actual coverage",
    "demo.value.dbAsset": "Database level (asset)",
    "demo.value.rpoRtoA": "< 2 min / 1.5 h",
    "demo.value.targetMysql": "Cross-host MySQL restore env",
    "demo.value.coverageTable": "order_details table only",
    "demo.value.rpoRtoB": "< 2 min / 4 h",
    "demo.value.targetProd": "Production OrderDB",
    "demo.value.coverageFullDb": "Entire OrderDB",
    "demo.value.rpoRtoC": "< 5 min / 6 h",
    "demo.value.targetIso": "Isolated validation environment",
    "demo.value.coverageValidate": "Full DB validate, then import target tables",
    "demo.layout.heading": "3 restore candidates generated for you",
    "demo.layout.intro":
      "This demo renders the updated layout-tree protocol without hard-coded candidate cards from business code.",
    "demo.layout.badgeRecommended": "Recommended",
    "demo.layout.badgeMediumRisk": "Medium risk",
    "demo.layout.paraA": "Fast restore with impact limited to the target business table.",
    "demo.metric.rpo": "RPO",
    "demo.metric.rto": "RTO",
    "demo.metric.success": "Success rate",
    "demo.metric.valueRpo": "< 2 min",
    "demo.metric.valueRto": "~ 1.5 hours",
    "demo.metric.valueHigh": "High",
    "demo.layout.kvFullDb": "Full database",
    "demo.layout.kvResourceHigh": "High",
    "demo.layout.kvComplexityMid": "Medium",
    "demo.layout.resourceUsageLabel": "Resource usage",
    "demo.layout.execComplexityLabel": "Execution complexity",
    "demo.layout.badgeHighResource": "High resource use",
    "demo.layout.paraB": "Stable path but broader scope and higher temporary resources.",
    "demo.layout.titleB": "Option B · Latest snapshot full restore then filter tables",
    "demo.layout.badgeConsistency": "Consistency risk",
    "demo.layout.paraC": "Lower implementation cost but stricter consistency checks.",
    "demo.layout.titleC": "Option C · Logical export backfill",
    "demo.layout.riskTitle": "Risk note",
    "demo.layout.riskBody": "Extra validation is required for import results and related data.",
    "demo.layout.calloutBizTitle": "Business impact",
    "demo.layout.calloutBizText": "Recommended option reduces impact on other production objects.",
    "demo.action.confirmA": "Confirm option A",
    "demo.action.viewDetail": "View details",
    "demo.action.confirmB": "Confirm option B",
    "demo.action.confirmC": "Confirm option C",
  },
  "zh-CN": {
    "common.appName": "Anybackup",
    "common.appTagline": "Always Resilient",
    "common.loading": "加载中...",
    "login.username": "用户名",
    "login.password": "密码",
    "login.usernamePlaceholder": "请输入用户名",
    "login.passwordPlaceholder": "请输入密码",
    "login.showPassword": "显示密码",
    "login.hidePassword": "隐藏密码",
    "login.submit": "登录",
    "login.submitting": "登录中",
    "login.serviceUnavailable": "登录服务暂时不可用，请稍后重试",
    "validation.requiredUsername": "请输入用户名",
    "validation.requiredPassword": "请输入密码",
    "validation.passwordMinLength": "密码长度不能少于 8 位",
    "validation.passwordSameAsUsername": "密码不能与用户名完全相同",
    "validation.passwordNeedsLetterAndDigit": "密码需至少包含字母和数字",
    "validation.passwordConfirmationMismatch": "两次输入的密码不一致",
    "sidebar.demoUser": "演示用户",
    "sidebar.backupAdmin": "备份管理员",
    "sidebar.loginAccount": "登录账号",
    "sidebar.logout": "退出登录",
    "sidebar.newConversation": "新建会话",
    "sidebar.searchConversation": "搜索会话",
    "sidebar.refreshingConversations": "正在刷新会话列表...",
    "sidebar.noMatchedConversation": "没有找到匹配的历史会话，换个关键词试试。",
    "sidebar.noConversationYet": "当前还没有正式会话，先从新建会话开始。",
    "sidebar.historyConversation": "历史会话",
    "sidebar.expandSidebar": "展开侧栏",
    "sidebar.collapseSidebar": "收起侧栏",
    "sidebar.settings": "设置",
    "settings.title": "设置",
    "settings.description": "管理系统基础能力和用户账号。",
    "settings.userManagementTitle": "用户管理",
    "settings.userManagementDescription": "查看用户、新建用户、启用或停用账号，并由管理员直接重置密码。",
    "guide.progressLabel": "引导进度",
    "guide.startLearn": "开始了解",
    "guide.howItWorksTitle": "工作流程",
    "guide.howItWorksDescription": "自然语言驱动，AI 分析生成，人类审核决策",
    "guide.viewDemo": "查看演示",
    "guide.userDescribeTitle": "用户描述需求",
    "guide.userDescribeDescription": "像和专家对话一样，自然语言描述你的灾备需求",
    "guide.agentGenerateTitle": "Agent 生成方案",
    "guide.agentGenerateDescription": "AI 分析环境，自动生成最优备份方案",
    "guide.pendingReview": "待审核",
    "guide.nextReview": "下一步：审核方案",
    "guide.agentGeneratingPlan": "Agent 正在分析环境并生成方案...",
    "guide.reviewTitle": "AI 生成，人类决策",
    "guide.reviewDescription": "所有关键操作都由你最终确认",
    "guide.reviewPassed": "审核通过",
    "guide.planApproved": "方案已通过",
    "guide.agentExecuteBackup": "Agent 将自动开始执行备份任务",
    "guide.finishGuide": "完成引导",
    "guide.readyTitle": "一切就绪！",
    "guide.readyDescription": "三栏融合工作台让你同时感知资产状态、与 Agent 对话、审核方案。",
    "guide.enterWorkspace": "进入工作台",
    "guide.dontShowAgain": "下次不再提示",
    "guide.skipGuide": "跳过引导，直接进入",
    "guide.workflowStep1Label": "自然语言描述需求",
    "guide.workflowStep1Sub": "AI 理解你的灾备意图",
    "guide.workflowStep2Label": "Agent 分析环境",
    "guide.workflowStep2Sub": "自动扫描、推荐最优策略",
    "guide.workflowStep3Label": "审核方案",
    "guide.workflowStep3Sub": "AI 生成、人类决策",
    "guide.workflowStep4Label": "自动执行与验证",
    "guide.workflowStep4Sub": "7x24 持续保护与巡检",
    "guide.product1Desc": "AI 原生自主备份恢复智能体",
    "guide.product2Desc": "专为智能体操作设计的备份恢复系统",
    "guide.product3Desc": "资产备份代理",
    "guide.feature1": "自主感知",
    "guide.feature2": "零门槛",
    "guide.feature3": "7x24 在线",
    "guide.feature4": "透明可信",
    "guide.feature5": "自主进化",
    "guide.userText": "ERP 的 Oracle 数据库需要做备份保护",
    "guide.agentText": "好的。我看到环境中有 ERP Oracle 19c RAC。\n\nRPO/RTO 要求是多少？可用的备份窗口是什么？",
    "guide.planTitle": "ERP Oracle RAC 备份方案",
    "guide.planEnvironment": "Oracle 19c RAC",
    "guide.planMetricBackupType": "备份类型",
    "guide.planMetricStrategy": "策略",
    "guide.planMetricRpoRto": "RPO / RTO",
    "guide.planMetricDataReduction": "数据缩减",
    "guide.planMetricBackupTypeValue": "RMAN Full+Incr+Log",
    "guide.planMetricStrategyValue": "每日 02:00 全量",
    "guide.planMetricRpoRtoValue": "<5min / <30min",
    "guide.planMetricDataReductionValue": "源端重删",
    "guide.planRecommendation": "此方案 RTO 优于同类系统均值 20%，推荐采用",
    "guide.reviewPlanSummary": "RMAN Full+Incr+Log · RPO<5min · 源端重删",
    "validation.requiredDisplayName": "请输入显示名称",
    "validation.confirmPasswordRequired": "请确认密码",
    "errors.genericUnavailable": "服务暂时不可用，请稍后重试。",
    "errors.networkUnavailable": "网络不可用，请检查网络后重试。",
    "conversation.placeholder.richContent": "[结构化内容]",
    "conversation.placeholder.clarification": "[待确认内容]",
    "conversation.placeholder.status": "[状态更新]",
    "conversation.error.listLoadFailed": "会话列表加载失败，请稍后重试",
    "conversation.error.detailLoadFailed": "会话详情加载失败，请稍后重试",
    "conversation.error.messagesLoadFailed": "会话消息加载失败，请稍后重试",
    "conversation.error.eventsLoadFailed": "会话状态事件加载失败，请稍后重试",
    "conversation.error.createFailed": "会话创建失败，请稍后重试",
    "conversation.error.sendFailed": "消息发送失败，请稍后重试",
    "conversation.error.archiveFailed": "会话归档失败，请稍后重试",
    "conversation.error.restoreFailed": "会话恢复失败，请稍后重试",
    "conversation.error.copyConfigFailed": "会话配置复制失败，请稍后重试",
    "conversation.pendingTurnTimeout": "响应超时，请稍后重试。",
    "conversation.serviceUnavailable": "会话服务暂时不可用，请稍后重试。",
    "conversation.newConversationTitle": "新对话",
    "auth.authServiceUnavailable": "认证服务暂时不可用，请稍后重试",
    "auth.userDisabled": "账号不可用，请联系管理员",
    "auth.invalidCredentials": "用户名或密码不正确",
    "auth.sessionValidationFailed": "登录态校验失败，请重新登录",
    "auth.sessionExpired": "登录已失效，请重新登录。",
    "auth.logoutFailed": "退出登录失败，请稍后重试",
    "users.unnamed": "未命名用户",
    "users.listLoadFailed": "用户列表加载失败，请稍后重试",
    "users.detailLoadFailed": "用户信息加载失败，请稍后重试",
    "users.roleLoadFailed": "内置管理员角色详情加载失败，请稍后重试",
    "users.adminRoleNotFound":
      "未找到可分配的内置管理员角色（已尝试：{{roles}}）。请检查当前环境的角色配置。",
    "users.assignRoleFailed": "内置管理员角色分配失败，请稍后重试",
    "users.createReadFailed": "用户创建成功，但未能读取用户详情",
    "users.createFailed": "用户创建失败，请稍后重试",
    "users.assignAfterCreateFailed": "用户已创建，但分配内置管理员角色失败，请稍后重试",
    "users.usernameExists": "用户名已存在",
    "users.saveFailed": "用户信息保存失败，请稍后重试",
    "users.statusSaveFailed": "用户状态保存失败，请稍后重试",
    "users.cannotDisableSelf": "不能停用当前登录用户",
    "users.resetPasswordFailed": "密码重置失败，请稍后重试",
    "users.feedback.created": "用户已创建。",
    "users.feedback.saved": "用户信息已保存。",
    "users.feedback.enabled": "用户已启用。",
    "users.feedback.disabled": "用户已停用。",
    "users.feedback.passwordReset": "密码已重置，请通知用户使用新密码登录。",
    "chat.closeHint": "关闭提示",
    "chat.preparingWorkspace": "正在准备会话工作台...",
    "chat.restoringConversation": "正在恢复历史会话...",
    "chat.refNotInMessagePrefix": "当前引用内容未随消息返回，暂时无法打开：",
    "chat.refreshingList": "正在刷新会话列表...",
    "chat.conversationNotLoaded": "当前会话不存在或还未加载完成，请重新选择。",
    "chat.composer.placeholder": "告诉我你想备份、恢复或查询什么...",
    "chat.composer.hint": "基于当前会话继续提问，候选方案确认也会作为结构化输入继续推进。",
    "chat.composer.send": "发送消息",
    "chat.composer.inputLabel": "消息输入",
    "chat.empty.title": "从这里开始一段新对话",
    "chat.empty.description": "先从当前会话继续提问，后续候选方案确认也会在同一轮对话里推进。",
    "chat.panelHeader.title": "对话入口",
    "chat.panelHeader.subtitle": "对话逻辑将在下一个 plan 单独设计。",
    "chat.waiting.thinking": "思考中",
    "chat.waiting.duration": "已思考 {minutes}m{seconds}s",
    "chat.thought.label": "思考",
    "chat.thought.running": "思考中",
    "chat.thought.done": "已完成",
    "chat.messageTime.invalid": "刚刚",
    "candidate.recommended": "推荐",
    "candidate.confirmed": "已确认",
    "workspace.breadcrumb.home": "工作台",
    "workspace.breadcrumb.settings": "设置",
    "workspace.breadcrumb.users": "用户管理",
    "workspace.navTitle": "首页",
    "modal.close": "关闭",
    "language.zhShort": "中文",
    "usersPage.title": "用户管理",
    "usersPage.subtitle":
      "管理登录账号、初始化密码和启停状态。新建用户后会自动分配管理员角色。",
    "usersPage.enabledCount": "已启用 {enabled} / {total}",
    "usersPage.createUser": "新建用户",
    "usersPage.listTitle": "用户列表",
    "usersPage.listHint": "当前页面只展示登录账号管理相关信息。",
    "usersPage.userCount": "共 {count} 位用户",
    "usersPage.colUsername": "用户名",
    "usersPage.colDisplayName": "显示名称",
    "usersPage.colRole": "角色",
    "usersPage.colStatus": "状态",
    "usersPage.colActions": "操作",
    "usersPage.loading": "正在加载用户列表...",
    "usersPage.empty": "暂无用户",
    "usersPage.currentAccount": "当前账号",
    "usersPage.edit": "编辑",
    "usersPage.disable": "停用",
    "usersPage.enable": "启用",
    "usersPage.resetPassword": "重置密码",
    "userForm.titleCreate": "新建用户",
    "userForm.titleEdit": "编辑用户",
    "userForm.descCreate": "管理员直接设置初始密码，创建后自动分配管理员角色。",
    "userForm.descEdit": "用户创建后用户名不可修改。",
    "userForm.cancel": "取消",
    "userForm.create": "创建用户",
    "userForm.save": "保存",
    "userForm.labelUsername": "用户名",
    "userForm.labelDisplayName": "显示名称",
    "userForm.labelPassword": "密码",
    "userForm.labelConfirmPassword": "确认密码",
    "userForm.status": "状态",
    "userForm.statusHint": "控制该用户是否允许继续登录系统。",
    "userForm.statusEnabled": "启用",
    "userForm.statusDisabled": "停用",
    "userForm.enableToDisableHint": "当前状态为启用，点击切换为停用",
    "userForm.disableToEnableHint": "当前状态为停用，点击切换为启用",
    "userForm.formError": "服务暂时不可用，请稍后重试。",
    "resetPassword.title": "重置密码",
    "resetPassword.newPassword": "新密码",
    "resetPassword.confirmNewPassword": "确认新密码",
    "resetPassword.description": "请为 {{username}} 设置新密码。系统不会自动发送密码。",
    "resetPassword.cancel": "取消",
    "resetPassword.submit": "重置密码",
    "resetPassword.targetLabel": "重置对象",
    "disableUser.title": "停用用户",
    "disableUser.description": "停用后该用户将无法登录系统，历史会话和账号资料仍会保留。",
    "disableUser.cancel": "取消",
    "disableUser.confirm": "确认停用",
    "disableUser.warningTitle": "停用后将立即失去后台访问权限。",
    "disableUser.warningBody": "后续如果需要恢复使用，可以在用户管理中重新启用该账号。",
    "disableUser.preview": "即将停用用户 {{username}}",
    "disableUser.selectUser": "请选择要停用的用户。",
    "demo.pageTitle": "对话状态演示 Demo",
    "demo.pageSubtitle":
      "这个页面不依赖真实会话接口，直接复用现有的聊天渲染组件，方便快速查看普通消息、等待态和候选方案三种形态的实际样子。",
    "demo.visitPath": "访问路径：",
    "demo.sectionLayoutTreeDesc": "基于更新后的协议，用固定节点注册表和动作分发器渲染候选方案对比块。",
    "demo.protocolActionFeedback": "协议动作反馈：",
    "demo.sectionLayoutTitle": "AG-UI Layout Tree Demo",
    "demo.sectionNodeGalleryDesc":
      "把协议白名单里剩余的节点类型集中展示出来，方便直接看 tabs、markdown、table、chart、attachment 和动作类型。",
    "demo.sectionCandidateTitle": "形态一：候选方案确认卡片",
    "demo.sectionCandidateDesc":
      "这是结构化富内容的主场景。当前页面会本地接住确认、拒绝和补充约束动作，模拟真实交互闭环。",
    "demo.lastInteraction": "最近一次交互：",
    "demo.sectionTextTitle": "形态二：普通消息",
    "demo.sectionTextDesc": "用于最常见的一问一答文本回复，展示基础气泡、标题栏和消息时间。",
    "demo.sectionWaitingTitle": "形态三：等待回复",
    "demo.sectionWaitingDesc":
      "用户消息已经发出，但 assistant 还在处理中，底部会显示等待中的占位反馈。",
    "demo.listTitleRestore": "订单数据库恢复",
    "demo.listSummaryCandidate": "候选方案消息演示",
    "demo.listTitleBackupQuery": "备份状态查询",
    "demo.listSummaryPlain": "普通文本消息演示",
    "demo.listTitleGenerating": "恢复方案生成中",
    "demo.listSummaryWaiting": "消息发送后，等待模型继续回复",
    "demo.selection.none": "还没有触发交互，你可以点候选卡片里的按钮看看真实反馈。",
    "demo.selection.confirm": "已模拟提交确认：",
    "demo.selection.reject": "已模拟提交拒绝：",
    "demo.selection.revise": "已模拟提交补充约束：",
    "demo.selection.reviseEmpty": "未填写",
    "demo.layout.none": "还没有触发动作。你可以直接点卡片里的确认按钮看协议动作怎么落到前端。",
    "demo.layout.last": "最近一次动作：",
    "demo.text.user.backupStatus": "帮我看看今天凌晨订单库的备份状态。",
    "demo.text.assistant.backupStatus":
      "今天凌晨 02:00 的订单库备份已成功完成，耗时 7 分 12 秒。当前没有发现异常，可以继续查看恢复点或导出详情。",
    "demo.waiting.user.restorePlan": "请帮我生成恢复订单数据库到昨天下午的方案。",
    "demo.candidate.user.message": "恢复到 2026-04-19 下午，先给我候选方案。",
    "demo.candidate.assistant.lead": "我找到了 3 个可用恢复点。",
    "demo.candidate.cardTitle": "为你生成 3 个备选方案",
    "demo.candidate.cardSummary": "已根据恢复粒度、业务影响和执行成本筛出 3 个候选项。",
    "demo.candidate.actionConfirm": "确认提交方案",
    "demo.candidate.actionReject": "放弃该方案",
    "demo.candidate.actionRevise": "补充约束",
    "demo.candidate.reviseInputLabel": "补充约束",
    "demo.candidate.revisePlaceholder": "例如：先生成方案，不要立即执行。",
    "demo.candidate.reviseSubmit": "提交补充约束",
    "demo.optionA.title": "方案 A：异机数据库级恢复 + 表导出导入",
    "demo.optionA.summary": "推荐方案，既能控制恢复范围，也能缩短业务窗口。",
    "demo.optionA.extraTitle": "业务影响",
    "demo.optionA.extraBody":
      "恢复粒度为数据库级，但实际仅操作目标表，能避开生产库其他表数据安全风险。",
    "demo.optionB.title": "方案 B：直接数据库级恢复到生产库",
    "demo.optionB.summary": "恢复路径最短，但会直接影响线上库。",
    "demo.optionB.extraTitle": "业务影响",
    "demo.optionB.extraBody": "直接恢复会覆盖其他正常业务表数据，风险较高。",
    "demo.optionC.title": "方案 C：延迟窗口内执行异机校验恢复",
    "demo.optionC.summary": "适合需要先验证数据完整性的场景。",
    "demo.optionC.extraTitle": "业务影响",
    "demo.optionC.extraBody": "整体耗时更长，但能在正式恢复前完成一致性校验。",
    "demo.field.restoreScope": "恢复粒度",
    "demo.field.rpoRto": "RPO / RTO",
    "demo.field.target": "恢复目的地",
    "demo.field.coverage": "实际恢复范围",
    "demo.value.dbAsset": "数据库级（资产）",
    "demo.value.rpoRtoA": "< 2 分钟 / 1.5 小时",
    "demo.value.targetMysql": "异机 MySQL 恢复环境",
    "demo.value.coverageTable": "仅 order_details 表",
    "demo.value.rpoRtoB": "< 2 分钟 / 4 小时",
    "demo.value.targetProd": "生产 OrderDB",
    "demo.value.coverageFullDb": "整个 OrderDB",
    "demo.value.rpoRtoC": "< 5 分钟 / 6 小时",
    "demo.value.targetIso": "隔离校验环境",
    "demo.value.coverageValidate": "全库校验后再导回目标表",
    "demo.layout.heading": "已为你生成 3 个恢复候选方案",
    "demo.layout.intro":
      "这个 demo 直接按照更新后的 layout-tree 协议渲染，不再依赖业务写死的候选方案卡。",
    "demo.layout.badgeRecommended": "推荐",
    "demo.layout.badgeMediumRisk": "中风险",
    "demo.layout.paraA": "恢复速度较快，且可将影响范围收敛到目标业务表。",
    "demo.metric.rpo": "RPO",
    "demo.metric.rto": "RTO",
    "demo.metric.success": "成功率",
    "demo.metric.valueRpo": "< 2 分钟",
    "demo.metric.valueRto": "约 1.5 小时",
    "demo.metric.valueHigh": "高",
    "demo.layout.kvFullDb": "整库",
    "demo.layout.kvResourceHigh": "高",
    "demo.layout.kvComplexityMid": "中",
    "demo.layout.resourceUsageLabel": "资源占用",
    "demo.layout.execComplexityLabel": "执行复杂度",
    "demo.layout.badgeHighResource": "高资源消耗",
    "demo.layout.paraB": "实现路径稳定，但恢复范围较大，对临时资源要求更高。",
    "demo.layout.titleB": "方案 B：最近快照整库恢复后再筛表",
    "demo.layout.badgeConsistency": "一致性风险",
    "demo.layout.paraC": "实施成本低，但对数据一致性和校验要求更高。",
    "demo.layout.titleC": "方案 C：逻辑导出补数",
    "demo.layout.riskTitle": "风险提示",
    "demo.layout.riskBody": "需要额外校验导入结果与关联数据完整性。",
    "demo.layout.calloutBizTitle": "业务影响",
    "demo.layout.calloutBizText": "推荐方案，能降低对生产库其他对象的影响。",
    "demo.action.confirmA": "确认方案 A",
    "demo.action.viewDetail": "查看详细说明",
    "demo.action.confirmB": "确认方案 B",
    "demo.action.confirmC": "确认方案 C",
  },
}

export function translate(key: MessageKey, locale: Locale = getStoredLocale()): string {
  return messages[locale][key] ?? messages.en[key]
}

