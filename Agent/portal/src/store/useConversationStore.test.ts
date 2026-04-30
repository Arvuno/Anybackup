import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"

const {
  createConversationMock,
  getConversationDetailMock,
  getConversationMessagesMock,
  listConversationEventsMock,
  listConversationsMock,
  searchConversationsMock,
  sendMessageMock,
} = vi.hoisted(() => ({
  createConversationMock: vi.fn(),
  getConversationDetailMock: vi.fn(),
  getConversationMessagesMock: vi.fn(),
  listConversationEventsMock: vi.fn(),
  listConversationsMock: vi.fn(),
  searchConversationsMock: vi.fn(),
  sendMessageMock: vi.fn(),
}))

vi.mock("@/services/conversation-service", () => ({
  createConversation: createConversationMock,
  getConversationDetail: getConversationDetailMock,
  getConversationMessages: getConversationMessagesMock,
  listConversationEvents: listConversationEventsMock,
  listConversations: listConversationsMock,
  searchConversations: searchConversationsMock,
  sendMessage: sendMessageMock,
}))

import { conversationDraftKeyForConversation } from "@/lib/conversation-draft"
import { useConversationStore } from "@/store/useConversationStore"

function resetConversationStore(): void {
  useConversationStore.setState({
    bootstrapped: false,
    listLoading: false,
    conversationLoading: false,
    query: "",
    error: null,
    conversations: [],
    selectedWorkspace: null,
    localDraftWorkspace: null,
    detailsByConversationId: {},
    messagesByConversationId: {},
    draftsByKey: {},
    pendingTurnByConversationId: {},
    pendingTurnStartedAtMsByConversationId: {},
    nextPollAfterMsByConversationId: {},
    eventCursorByConversationId: {},
    latestEventSequenceByConversationId: {},
    appliedStatusEventIdsByConversationId: {},
    submittingWorkspaceKey: null,
  })
}

async function waitForPollCycle(): Promise<void> {
  await new Promise((resolve) => setTimeout(resolve, 0))
}

function messageCreatedStatusEvent(
  conversationId: string,
  turnId: string,
  messageId: string,
  sequence = 1,
) {
  return {
    statusEventId: `evt_${messageId}`,
    conversationId,
    turnId,
    messageId,
    eventType: "message.created" as const,
    sequence,
    interactionState: "thinking" as const,
    messageStatus: "published" as const,
    createdAt: "2026-04-23T09:00:00.000Z",
  }
}

describe("useConversationStore", () => {
  beforeEach(() => {
    localStorage.clear()
    resetConversationStore()

    createConversationMock.mockReset()
    getConversationDetailMock.mockReset()
    getConversationMessagesMock.mockReset()
    listConversationEventsMock.mockReset()
    listConversationsMock.mockReset()
    searchConversationsMock.mockReset()
    sendMessageMock.mockReset()

    listConversationsMock.mockResolvedValue([])
    searchConversationsMock.mockResolvedValue([])
    getConversationDetailMock.mockResolvedValue({
      conversationId: "conv_restore_001",
      title: "订单数据库恢复",
      createdAt: "2026-04-22T10:00:00.000Z",
      updatedAt: "2026-04-22T10:00:00.000Z",
      interactionState: "idle",
    })
    getConversationMessagesMock.mockResolvedValue([])
    listConversationEventsMock.mockResolvedValue({
      events: [],
      nextCursor: null,
      hasMore: false,
      latestSequence: 0,
      recommendedPollIntervalMs: 1000,
      interactionState: "completed",
    })
  })

  afterEach(async () => {
    if (useConversationStore.getState().selectedWorkspace?.kind === "conversation") {
      useConversationStore.getState().activateLocalDraftWorkspace()
    }
    await waitForPollCycle()
    vi.useRealTimers()
    vi.clearAllMocks()
    localStorage.clear()
    resetConversationStore()
  })

  it("promotes a local draft into a formal conversation on first send", async () => {
    createConversationMock.mockResolvedValue({
      conversation: {
        conversationId: "conv_restore_001",
        title: "订单数据库恢复",
        createdAt: "2026-04-22T10:00:00.000Z",
        updatedAt: "2026-04-22T10:00:00.000Z",
        interactionState: "thinking",
        activeTurnId: "turn_001",
      },
      message: {
        messageId: "msg_user_001",
        conversationId: "conv_restore_001",
        turnId: "turn_001",
        role: "user",
        contentType: "text",
        content: "帮我恢复订单数据库",
        createdAt: "2026-04-22T10:00:00.000Z",
        status: "published",
      },
      statusEvent: messageCreatedStatusEvent("conv_restore_001", "turn_001", "msg_user_001"),
      nextPollAfterMs: 0,
    })

    await useConversationStore.getState().hydrate()
    useConversationStore.getState().setDraft("帮我恢复订单数据库")

    await useConversationStore.getState().submitComposerMessage()

    const state = useConversationStore.getState()
    expect(state.selectedWorkspace).toEqual({
      kind: "conversation",
      conversationId: "conv_restore_001",
    })
    expect(state.messagesByConversationId.conv_restore_001).toHaveLength(1)
    expect(state.messagesByConversationId.conv_restore_001?.[0]?.content).toBe("帮我恢复订单数据库")
    expect(state.pendingTurnByConversationId.conv_restore_001).toEqual({ state: "thinking", turnId: "turn_001" })
  })

  it("submits candidate selection back into the same conversation", async () => {
    sendMessageMock.mockResolvedValue({
      conversation: {
        conversationId: "conv_restore_001",
        title: "订单数据库恢复",
        createdAt: "2026-04-22T10:00:00.000Z",
        updatedAt: "2026-04-22T10:01:00.000Z",
        interactionState: "thinking",
        activeTurnId: "turn_selection_001",
      },
      message: {
        messageId: "msg_user_selection_001",
        conversationId: "conv_restore_001",
        turnId: "turn_selection_001",
        role: "user",
        contentType: "text",
        content: "确认推荐方案：恢复到 2026-04-19 14:45",
        createdAt: "2026-04-22T10:01:00.000Z",
        status: "published",
      },
      statusEvent: messageCreatedStatusEvent(
        "conv_restore_001",
        "turn_selection_001",
        "msg_user_selection_001",
      ),
      nextPollAfterMs: 0,
    })

    useConversationStore.setState({
      bootstrapped: true,
      selectedWorkspace: {
        kind: "conversation",
        conversationId: "conv_restore_001",
      },
      detailsByConversationId: {
        conv_restore_001: {
          conversationId: "conv_restore_001",
          title: "订单数据库恢复",
          createdAt: "2026-04-22T10:00:00.000Z",
          updatedAt: "2026-04-22T10:00:00.000Z",
        },
      },
      messagesByConversationId: {
        conv_restore_001: [
          {
            messageId: "msg_candidate_001",
            conversationId: "conv_restore_001",
            role: "assistant",
            contentType: "rich_content",
            content: "已生成 2 个恢复候选方案。",
            createdAt: "2026-04-22T10:00:00.000Z",
            richPayload: {
              kind: "candidate_options",
              data: {
                reasoningTraceId: "trace_restore_001",
                title: "恢复候选方案",
                summary: "已找到 2 个可用备份点",
                actions: [
                  { type: "confirm", label: "确认提交方案" },
                  { type: "reject", label: "放弃该方案" },
                  {
                    type: "revise",
                    label: "补充约束",
                    inputLabel: "补充约束",
                    inputPlaceholder: "例如：先生成方案，不要立即执行。",
                    submitLabel: "提交补充约束",
                  },
                ],
                options: [
                  {
                    optionId: "option_a",
                    title: "方案 A：异机数据库级恢复 + 表导出导入",
                    recommended: true,
                    summary: "推荐方案",
                    fields: [
                      {
                        key: "restore_scope",
                        label: "恢复粒度",
                        value: "数据库级（资产）",
                      },
                    ],
                    extra: {
                      title: "业务影响",
                      content: "恢复粒度为数据库级，但实际仅操作目标表。",
                    },
                  },
                ],
              },
            },
          },
        ],
      },
      pendingTurnByConversationId: {
        conv_restore_001: {
          state: "clarifying",
          sourceMessageId: "msg_candidate_001",
        },
      },
    })

    await useConversationStore.getState().submitCandidateSelection({
      type: "candidate_selection",
      messageId: "msg_candidate_001",
      reasoningTraceId: "trace_restore_001",
      candidateOptionId: "option_a",
      selection: "confirm",
    })

    expect(sendMessageMock).toHaveBeenCalledWith(
      "conv_restore_001",
      expect.objectContaining({
        type: "candidate_selection",
        candidateOptionId: "option_a",
        selection: "confirm",
      }),
    )

    const state = useConversationStore.getState()
    expect(state.messagesByConversationId.conv_restore_001).toHaveLength(2)
    expect(state.messagesByConversationId.conv_restore_001?.[1]?.content).toBe(
      "确认推荐方案：恢复到 2026-04-19 14:45",
    )
    expect(state.pendingTurnByConversationId.conv_restore_001).toEqual({
      state: "thinking",
      turnId: "turn_selection_001",
    })
  })

  it("submits clarification responses back into the same conversation", async () => {
    sendMessageMock.mockResolvedValue({
      conversation: {
        conversationId: "conv_restore_001",
        title: "恢复订单数据库",
        createdAt: "2026-04-22T10:00:00.000Z",
        updatedAt: "2026-04-22T10:02:00.000Z",
        interactionState: "thinking",
        activeTurnId: "turn_clarification_001",
      },
      message: {
        messageId: "msg_user_clarification_001",
        conversationId: "conv_restore_001",
        turnId: "turn_clarification_001",
        role: "user",
        contentType: "text",
        content: "确认恢复窗口：最近安全点",
        createdAt: "2026-04-22T10:02:00.000Z",
        status: "published",
      },
      statusEvent: messageCreatedStatusEvent(
        "conv_restore_001",
        "turn_clarification_001",
        "msg_user_clarification_001",
      ),
      nextPollAfterMs: 0,
    })

    useConversationStore.setState({
      bootstrapped: true,
      selectedWorkspace: {
        kind: "conversation",
        conversationId: "conv_restore_001",
      },
      detailsByConversationId: {
        conv_restore_001: {
          conversationId: "conv_restore_001",
          title: "恢复订单数据库",
          createdAt: "2026-04-22T10:00:00.000Z",
          updatedAt: "2026-04-22T10:00:00.000Z",
        },
      },
      messagesByConversationId: {
        conv_restore_001: [
          {
            messageId: "msg_clarification_prompt_001",
            conversationId: "conv_restore_001",
            role: "assistant",
            contentType: "clarification",
            content: "请确认目标恢复时间窗口。",
            createdAt: "2026-04-22T10:01:00.000Z",
            richPayload: {
              kind: "clarification",
              data: {
                clarificationId: "restore_window",
                prompt: "请确认目标恢复时间窗口。",
                options: [
                  { label: "最近安全点", value: "latest_safe_point" },
                  { label: "指定时间", value: "custom_timestamp" },
                ],
                inputConstraints: {
                  required: true,
                  allowFreeText: true,
                },
              },
            },
          },
        ],
      },
      pendingTurnByConversationId: {
        conv_restore_001: {
          state: "clarifying",
          sourceMessageId: "msg_clarification_prompt_001",
        },
      },
    })

    await useConversationStore.getState().submitClarificationResponse({
      type: "clarification_response",
      messageId: "msg_clarification_prompt_001",
      clarificationId: "restore_window",
      selectedValue: "latest_safe_point",
      freeText: "最近安全点",
    })

    expect(sendMessageMock).toHaveBeenCalledWith(
      "conv_restore_001",
      expect.objectContaining({
        type: "clarification_response",
        clarificationId: "restore_window",
        selectedValue: "latest_safe_point",
      }),
    )

    const state = useConversationStore.getState()
    expect(state.messagesByConversationId.conv_restore_001).toHaveLength(2)
    expect(state.messagesByConversationId.conv_restore_001?.[1]?.content).toBe("确认恢复窗口：最近安全点")
    expect(state.pendingTurnByConversationId.conv_restore_001).toEqual({
      state: "thinking",
      turnId: "turn_clarification_001",
    })
  })

  it("hydrates an assistant text reply after the first-send polling cycle", async () => {
    createConversationMock.mockResolvedValue({
      conversation: {
        conversationId: "conv_text_loop",
        title: "Text loop",
        createdAt: "2026-04-23T09:00:00.000Z",
        updatedAt: "2026-04-23T09:00:00.000Z",
        interactionState: "thinking",
        activeTurnId: "turn_text_001",
      },
      message: {
        messageId: "msg_user_001",
        conversationId: "conv_text_loop",
        turnId: "turn_text_001",
        role: "user",
        contentType: "text",
        content: "help me restore the order database",
        createdAt: "2026-04-23T09:00:00.000Z",
        status: "published",
      },
      statusEvent: messageCreatedStatusEvent("conv_text_loop", "turn_text_001", "msg_user_001"),
      nextPollAfterMs: 0,
    })

    listConversationEventsMock.mockResolvedValueOnce({
      events: [
        {
          statusEventId: "evt_assistant_001",
          conversationId: "conv_text_loop",
          turnId: "turn_text_001",
          messageId: "msg_assistant_001",
          eventType: "message.created",
          sequence: 2,
          interactionState: "thinking",
          messageStatus: "responded",
          message: {
            messageId: "msg_assistant_001",
            conversationId: "conv_text_loop",
            turnId: "turn_text_001",
            role: "assistant",
            contentType: "text",
            content: "I will first inspect the latest restore points.",
            createdAt: "2026-04-23T09:00:20.000Z",
            status: "responded",
          },
          createdAt: "2026-04-23T09:00:20.000Z",
        },
        {
          statusEventId: "evt_completed_001",
          conversationId: "conv_text_loop",
          turnId: "turn_text_001",
          eventType: "interaction.status_changed",
          sequence: 3,
          interactionState: "completed",
          activeTurnId: null,
          completedTurnId: "turn_text_001",
          createdAt: "2026-04-23T09:00:21.000Z",
        },
      ],
      nextCursor: "cursor_text_003",
      hasMore: false,
      latestSequence: 3,
      recommendedPollIntervalMs: 1000,
      interactionState: "completed",
    })

    await useConversationStore.getState().hydrate()
    useConversationStore.getState().setDraft("help me restore the order database")

    await useConversationStore.getState().submitComposerMessage()
    await waitForPollCycle()

    const state = useConversationStore.getState()
    expect(state.messagesByConversationId.conv_text_loop?.map((message) => message.content)).toEqual([
      "help me restore the order database",
      "I will first inspect the latest restore points.",
    ])
    expect(state.pendingTurnByConversationId.conv_text_loop).toEqual({ state: "idle" })
  })

  it("caps an accepted next poll delay at five seconds", async () => {
    vi.useFakeTimers()

    try {
      createConversationMock.mockResolvedValue({
        conversation: {
          conversationId: "conv_capped_accept_delay",
          title: "Capped accept delay",
          createdAt: "2026-04-23T09:00:00.000Z",
          updatedAt: "2026-04-23T09:00:00.000Z",
          interactionState: "thinking",
          activeTurnId: "turn_capped_accept_delay_001",
        },
        message: {
          messageId: "msg_user_capped_accept_delay_001",
          conversationId: "conv_capped_accept_delay",
          turnId: "turn_capped_accept_delay_001",
          role: "user",
          contentType: "text",
          content: "keep polling this restore flow",
          createdAt: "2026-04-23T09:00:00.000Z",
          status: "published",
        },
        statusEvent: messageCreatedStatusEvent(
          "conv_capped_accept_delay",
          "turn_capped_accept_delay_001",
          "msg_user_capped_accept_delay_001",
        ),
        nextPollAfterMs: 60000,
      })

      listConversationEventsMock.mockResolvedValueOnce({
        events: [],
        nextCursor: null,
        hasMore: false,
        latestSequence: 1,
        recommendedPollIntervalMs: 1000,
        interactionState: "completed",
      })

      await useConversationStore.getState().hydrate()
      useConversationStore.getState().setDraft("keep polling this restore flow")

      await useConversationStore.getState().submitComposerMessage()

      expect(listConversationEventsMock).toHaveBeenCalledTimes(0)

      await vi.advanceTimersByTimeAsync(4999)
      expect(listConversationEventsMock).toHaveBeenCalledTimes(0)

      await vi.advanceTimersByTimeAsync(1)
      expect(listConversationEventsMock).toHaveBeenCalledTimes(1)
    } finally {
      vi.useRealTimers()
    }
  })

  it("keeps the accepted user message before later assistant events even when event timestamps are earlier", async () => {
    createConversationMock.mockResolvedValue({
      conversation: {
        conversationId: "conv_timestamp_skew",
        title: "Timestamp skew",
        createdAt: "2026-04-23T09:00:00.000Z",
        updatedAt: "2026-04-23T09:00:10.000Z",
        interactionState: "thinking",
        activeTurnId: "turn_skew_001",
      },
      message: {
        messageId: "msg_user_skew_001",
        conversationId: "conv_timestamp_skew",
        turnId: "turn_skew_001",
        role: "user",
        contentType: "text",
        content: "66",
        createdAt: "2026-04-23T09:00:10.000Z",
        status: "published",
      },
      statusEvent: messageCreatedStatusEvent("conv_timestamp_skew", "turn_skew_001", "msg_user_skew_001"),
      nextPollAfterMs: 0,
    })

    listConversationEventsMock.mockResolvedValueOnce({
      events: [
        {
          statusEventId: "evt_assistant_skew_001",
          conversationId: "conv_timestamp_skew",
          turnId: "turn_skew_001",
          messageId: "msg_assistant_skew_001",
          eventType: "message.created",
          sequence: 2,
          interactionState: "thinking",
          messageStatus: "responded",
          message: {
            messageId: "msg_assistant_skew_001",
            conversationId: "conv_timestamp_skew",
            turnId: "turn_skew_001",
            role: "assistant",
            contentType: "text",
            content: "正在理解问题并准备查询上下文。",
            createdAt: "2026-04-23T09:00:00.000Z",
            status: "responded",
          },
          createdAt: "2026-04-23T09:00:00.000Z",
        },
        {
          statusEventId: "evt_completed_skew_001",
          conversationId: "conv_timestamp_skew",
          turnId: "turn_skew_001",
          eventType: "interaction.status_changed",
          sequence: 3,
          interactionState: "completed",
          activeTurnId: null,
          completedTurnId: "turn_skew_001",
          createdAt: "2026-04-23T09:00:11.000Z",
        },
      ],
      nextCursor: "cursor_skew_003",
      hasMore: false,
      latestSequence: 3,
      recommendedPollIntervalMs: 1000,
      interactionState: "completed",
    })

    await useConversationStore.getState().hydrate()
    useConversationStore.getState().setDraft("66")

    await useConversationStore.getState().submitComposerMessage()
    await waitForPollCycle()

    const state = useConversationStore.getState()
    expect(state.messagesByConversationId.conv_timestamp_skew?.map((message) => message.messageId)).toEqual([
      "msg_user_skew_001",
      "msg_assistant_skew_001",
    ])
  })

  it(
    "keeps polling after result content until an interaction completed event arrives",
    async () => {
    createConversationMock.mockResolvedValue({
      conversation: {
        conversationId: "conv_event_loop",
        title: "Event loop",
        createdAt: "2026-04-23T09:00:00.000Z",
        updatedAt: "2026-04-23T09:00:00.000Z",
        interactionState: "thinking",
        activeTurnId: "turn_001",
      },
      message: {
        messageId: "msg_user_001",
        conversationId: "conv_event_loop",
        turnId: "turn_001",
        role: "user",
        contentType: "text",
        content: "analyze restore options",
        createdAt: "2026-04-23T09:00:00.000Z",
        status: "published",
      },
      statusEvent: messageCreatedStatusEvent("conv_event_loop", "turn_001", "msg_user_001"),
      nextPollAfterMs: 0,
    })

    listConversationEventsMock
      .mockResolvedValueOnce({
        events: [
          {
            statusEventId: "evt_thought_001",
            conversationId: "conv_event_loop",
            turnId: "turn_001",
            messageId: "msg_thought_001",
            eventType: "rich_content.created",
            sequence: 2,
            interactionState: "thinking",
            messageStatus: "responded",
            message: {
              messageId: "msg_thought_001",
              conversationId: "conv_event_loop",
              turnId: "turn_001",
              parentMessageId: "msg_user_001",
              role: "assistant",
              contentType: "rich_content",
              content: "Checking restore point freshness...",
              createdAt: "2026-04-23T09:00:10.000Z",
              status: "responded",
              richPayload: {
                kind: "thought",
                data: {
                  summary: "Checking restore point freshness...",
                },
              },
            },
            createdAt: "2026-04-23T09:00:10.000Z",
          },
          {
            statusEventId: "evt_result_001",
            conversationId: "conv_event_loop",
            turnId: "turn_001",
            messageId: "msg_result_001",
            eventType: "rich_content.created",
            sequence: 3,
            interactionState: "executing",
            messageStatus: "responded",
            message: {
              messageId: "msg_result_001",
              conversationId: "conv_event_loop",
              turnId: "turn_001",
              parentMessageId: "msg_user_001",
              role: "assistant",
              contentType: "text",
              content: "The safest restore point is 2026-04-19 14:45.",
              createdAt: "2026-04-23T09:00:20.000Z",
              status: "responded",
            },
            createdAt: "2026-04-23T09:00:20.000Z",
          },
        ],
        nextCursor: "cursor_after_3",
        hasMore: false,
        latestSequence: 3,
        recommendedPollIntervalMs: 0,
        interactionState: "executing",
      })
      .mockResolvedValueOnce({
        events: [
          {
            statusEventId: "evt_completed_001",
            conversationId: "conv_event_loop",
            turnId: "turn_001",
            eventType: "interaction.status_changed",
            sequence: 4,
            interactionState: "completed",
            activeTurnId: null,
            completedTurnId: "turn_001",
            createdAt: "2026-04-23T09:00:21.000Z",
          },
        ],
        nextCursor: "cursor_after_4",
        hasMore: false,
        latestSequence: 4,
        recommendedPollIntervalMs: 1000,
        interactionState: "completed",
      })

    await useConversationStore.getState().hydrate()
    useConversationStore.getState().setDraft("analyze restore options")

    await useConversationStore.getState().submitComposerMessage()
    await waitForPollCycle()

    let state = useConversationStore.getState()
    expect(state.messagesByConversationId.conv_event_loop?.map((message) => message.messageId)).toEqual([
      "msg_user_001",
      "msg_thought_001",
      "msg_result_001",
    ])
    expect(state.pendingTurnByConversationId.conv_event_loop).toEqual({
      state: "thinking",
      turnId: "turn_001",
    })
    expect(listConversationEventsMock).toHaveBeenCalledTimes(1)

    await new Promise((resolve) => setTimeout(resolve, 25))
    expect(listConversationEventsMock).toHaveBeenCalledTimes(1)

    await new Promise((resolve) => setTimeout(resolve, 5100))

    state = useConversationStore.getState()
    expect(listConversationEventsMock).toHaveBeenCalledTimes(2)
    expect(state.pendingTurnByConversationId.conv_event_loop).toEqual({ state: "idle" })
    expect(state.eventCursorByConversationId.conv_event_loop).toBe("cursor_after_4")
    expect(state.latestEventSequenceByConversationId.conv_event_loop).toBe(4)
  },
  15_000,
  )

  it("caps long event repoll delays at five seconds and reports a timeout after fifteen minutes", async () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date("2026-04-26T02:00:00.000Z"))

    try {
      createConversationMock.mockResolvedValue({
        conversation: {
          conversationId: "conv_timeout",
          title: "Timeout loop",
          createdAt: "2026-04-26T02:00:00.000Z",
          updatedAt: "2026-04-26T02:00:00.000Z",
          interactionState: "thinking",
          activeTurnId: "turn_timeout_001",
        },
        message: {
          messageId: "msg_timeout_user_001",
          conversationId: "conv_timeout",
          turnId: "turn_timeout_001",
          role: "user",
          contentType: "text",
          content: "run long restore analysis",
          createdAt: "2026-04-26T02:00:00.000Z",
          status: "published",
        },
        statusEvent: messageCreatedStatusEvent("conv_timeout", "turn_timeout_001", "msg_timeout_user_001"),
        nextPollAfterMs: 0,
      })

      listConversationEventsMock.mockResolvedValue({
        events: [],
        nextCursor: null,
        hasMore: false,
        latestSequence: 1,
        recommendedPollIntervalMs: 120000,
        interactionState: "executing",
      })
      getConversationDetailMock.mockResolvedValue({
        conversationId: "conv_timeout",
        title: "Timeout loop",
        createdAt: "2026-04-26T02:00:00.000Z",
        updatedAt: "2026-04-26T02:00:00.000Z",
        interactionState: "executing",
        activeTurnId: "turn_timeout_001",
      })

      await useConversationStore.getState().hydrate()
      useConversationStore.getState().setDraft("run long restore analysis")

      await useConversationStore.getState().submitComposerMessage()
      await vi.advanceTimersByTimeAsync(0)

      let state = useConversationStore.getState()
      expect(listConversationEventsMock).toHaveBeenCalledTimes(1)
      expect(state.pendingTurnByConversationId.conv_timeout).toMatchObject({
        state: "thinking",
        turnId: "turn_timeout_001",
      })
      expect(
        (state as unknown as { pendingTurnStartedAtMsByConversationId: Record<string, number> })
          .pendingTurnStartedAtMsByConversationId.conv_timeout,
      ).toBe(new Date("2026-04-26T02:00:00.000Z").getTime())

      await vi.advanceTimersByTimeAsync(4999)
      expect(listConversationEventsMock).toHaveBeenCalledTimes(1)
      expect(useConversationStore.getState().pendingTurnByConversationId.conv_timeout).toMatchObject({
        state: "thinking",
      })

      await vi.advanceTimersByTimeAsync(1)
      expect(listConversationEventsMock).toHaveBeenCalledTimes(2)
      expect(useConversationStore.getState().pendingTurnByConversationId.conv_timeout).toMatchObject({
        state: "thinking",
      })

      await vi.advanceTimersByTimeAsync(894999)
      expect(listConversationEventsMock).toHaveBeenCalledTimes(180)
      expect(useConversationStore.getState().pendingTurnByConversationId.conv_timeout).toMatchObject({
        state: "thinking",
      })

      await vi.advanceTimersByTimeAsync(1)

      state = useConversationStore.getState()
      expect(listConversationEventsMock).toHaveBeenCalledTimes(180)
      expect(state.pendingTurnByConversationId.conv_timeout).toMatchObject({
        state: "error",
        turnId: "turn_timeout_001",
      })
      expect(state.error).toBe("The request timed out. Please try again later.")

      await vi.advanceTimersByTimeAsync(120000)
      expect(listConversationEventsMock).toHaveBeenCalledTimes(180)
    } finally {
      vi.useRealTimers()
    }
  })

  it("stops polling when a layout-tree state snapshot marks the current turn completed", async () => {
    createConversationMock.mockResolvedValue({
      conversation: {
        conversationId: "conv_layout_tree_state",
        title: "Layout tree state",
        createdAt: "2026-04-23T09:00:00.000Z",
        updatedAt: "2026-04-23T09:00:00.000Z",
        interactionState: "thinking",
        activeTurnId: "turn_layout_tree_state_001",
      },
      message: {
        messageId: "msg_user_layout_tree_state_001",
        conversationId: "conv_layout_tree_state",
        turnId: "turn_layout_tree_state_001",
        role: "user",
        contentType: "text",
        content: "show me the final restore recommendation",
        createdAt: "2026-04-23T09:00:00.000Z",
        status: "published",
      },
      statusEvent: messageCreatedStatusEvent(
        "conv_layout_tree_state",
        "turn_layout_tree_state_001",
        "msg_user_layout_tree_state_001",
      ),
      nextPollAfterMs: 0,
    })

    listConversationEventsMock.mockResolvedValueOnce({
      events: [
        {
          statusEventId: "evt_layout_tree_state_001",
          conversationId: "conv_layout_tree_state",
          turnId: "turn_layout_tree_state_001",
          messageId: "msg_layout_tree_state_001",
          eventType: "rich_content.created",
          sequence: 2,
          interactionState: "executing",
          messageStatus: "responded",
          message: {
            messageId: "msg_layout_tree_state_001",
            conversationId: "conv_layout_tree_state",
            turnId: "turn_layout_tree_state_001",
            parentMessageId: "msg_user_layout_tree_state_001",
            role: "assistant",
            contentType: "rich_content",
            content: "The final recommendation is ready.",
            createdAt: "2026-04-23T09:00:20.000Z",
            status: "responded",
            richPayload: {
              kind: "layout_tree",
              data: {
                activity: {
                  contract: "conversation.ui.layout-tree@1",
                  blockId: "restore_result",
                  ui: {
                    id: "root",
                    type: "stack",
                  },
                  meta: {
                    intent: "result",
                  },
                },
                stateSnapshot: {
                  interaction: {
                    status: "completed",
                  },
                },
              },
            },
          },
          createdAt: "2026-04-23T09:00:20.000Z",
        },
      ],
      nextCursor: "cursor_layout_tree_state_002",
      hasMore: false,
      latestSequence: 2,
      recommendedPollIntervalMs: 1000,
      interactionState: "executing",
    })

    await useConversationStore.getState().hydrate()
    useConversationStore.getState().setDraft("show me the final restore recommendation")

    await useConversationStore.getState().submitComposerMessage()
    await waitForPollCycle()

    const state = useConversationStore.getState()
    expect(listConversationEventsMock).toHaveBeenCalledTimes(1)
    expect(state.pendingTurnByConversationId.conv_layout_tree_state).toEqual({ state: "idle" })
    expect(state.detailsByConversationId.conv_layout_tree_state?.interactionState).toBe("completed")
    expect(state.detailsByConversationId.conv_layout_tree_state?.activeTurnId).toBeUndefined()
  })

  it(
    "reconciles control state from conversation detail after active polling stalls without a closing status event",
    async () => {
    createConversationMock.mockResolvedValue({
      conversation: {
        conversationId: "conv_reconcile_detail",
        title: "Reconcile detail",
        createdAt: "2026-04-23T09:00:00.000Z",
        updatedAt: "2026-04-23T09:00:00.000Z",
        interactionState: "thinking",
        activeTurnId: "turn_reconcile_detail_001",
      },
      message: {
        messageId: "msg_user_reconcile_detail_001",
        conversationId: "conv_reconcile_detail",
        turnId: "turn_reconcile_detail_001",
        role: "user",
        contentType: "text",
        content: "text fallback only",
        createdAt: "2026-04-23T09:00:00.000Z",
        status: "published",
      },
      statusEvent: messageCreatedStatusEvent(
        "conv_reconcile_detail",
        "turn_reconcile_detail_001",
        "msg_user_reconcile_detail_001",
      ),
      nextPollAfterMs: 0,
    })

    getConversationDetailMock.mockResolvedValueOnce({
      conversationId: "conv_reconcile_detail",
      title: "Reconcile detail",
      createdAt: "2026-04-23T09:00:00.000Z",
      updatedAt: "2026-04-23T09:00:30.000Z",
      interactionState: "completed",
    })

    listConversationEventsMock
      .mockResolvedValueOnce({
        events: [
          {
            statusEventId: "evt_reconcile_detail_001",
            conversationId: "conv_reconcile_detail",
            turnId: "turn_reconcile_detail_001",
            messageId: "msg_assistant_reconcile_detail_001",
            eventType: "rich_content.created",
            sequence: 2,
            interactionState: "executing",
            messageStatus: "responded",
            message: {
              messageId: "msg_assistant_reconcile_detail_001",
              conversationId: "conv_reconcile_detail",
              turnId: "turn_reconcile_detail_001",
              parentMessageId: "msg_user_reconcile_detail_001",
              role: "assistant",
              contentType: "text",
              content: "This response intentionally contains only AG-UI text message events.",
              createdAt: "2026-04-23T09:00:20.000Z",
              status: "responded",
            },
            createdAt: "2026-04-23T09:00:20.000Z",
          },
        ],
        nextCursor: "cursor_reconcile_detail_002",
        hasMore: false,
        latestSequence: 2,
        recommendedPollIntervalMs: 1000,
        interactionState: "executing",
      })
      .mockResolvedValueOnce({
        events: [],
        nextCursor: "cursor_reconcile_detail_002",
        hasMore: false,
        latestSequence: 2,
        recommendedPollIntervalMs: 1000,
        interactionState: "executing",
      })

    await useConversationStore.getState().hydrate()
    useConversationStore.getState().setDraft("text fallback only")

    await useConversationStore.getState().submitComposerMessage()
    await waitForPollCycle()

    let state = useConversationStore.getState()
    expect(listConversationEventsMock).toHaveBeenCalledTimes(1)
    expect(getConversationDetailMock).toHaveBeenCalledTimes(0)
    expect(state.pendingTurnByConversationId.conv_reconcile_detail).toEqual({
      state: "thinking",
      turnId: "turn_reconcile_detail_001",
    })

    await new Promise((resolve) => setTimeout(resolve, 5100))

    state = useConversationStore.getState()
    expect(listConversationEventsMock).toHaveBeenCalledTimes(2)
    expect(getConversationDetailMock).toHaveBeenCalledTimes(1)
    expect(state.pendingTurnByConversationId.conv_reconcile_detail).toEqual({ state: "idle" })
  },
  15_000,
  )

  it("stops active polling after a clarification event when the turn is waiting for user input", async () => {
    createConversationMock.mockResolvedValue({
      conversation: {
        conversationId: "conv_clarifying_loop",
        title: "Clarifying loop",
        createdAt: "2026-04-23T09:00:00.000Z",
        updatedAt: "2026-04-23T09:00:00.000Z",
        interactionState: "thinking",
        activeTurnId: "turn_clarifying_001",
      },
      message: {
        messageId: "msg_user_clarifying_001",
        conversationId: "conv_clarifying_loop",
        turnId: "turn_clarifying_001",
        role: "user",
        contentType: "text",
        content: "please clarify restore window",
        createdAt: "2026-04-23T09:00:00.000Z",
        status: "published",
      },
      statusEvent: messageCreatedStatusEvent(
        "conv_clarifying_loop",
        "turn_clarifying_001",
        "msg_user_clarifying_001",
      ),
      nextPollAfterMs: 0,
    })

    listConversationEventsMock.mockResolvedValueOnce({
      events: [
        {
          statusEventId: "evt_clarification_001",
          conversationId: "conv_clarifying_loop",
          turnId: "turn_clarifying_001",
          messageId: "msg_clarification_001",
          eventType: "rich_content.created",
          sequence: 2,
          interactionState: "clarifying",
          messageStatus: "responded",
          message: {
            messageId: "msg_clarification_001",
            conversationId: "conv_clarifying_loop",
            turnId: "turn_clarifying_001",
            parentMessageId: "msg_user_clarifying_001",
            role: "assistant",
            contentType: "clarification",
            content: "请确认目标恢复时间窗口。",
            createdAt: "2026-04-23T09:00:10.000Z",
            status: "responded",
          },
          createdAt: "2026-04-23T09:00:10.000Z",
        },
      ],
      nextCursor: "cursor_after_clarification",
      hasMore: false,
      latestSequence: 2,
      recommendedPollIntervalMs: 1000,
      interactionState: "clarifying",
    })

    await useConversationStore.getState().hydrate()
    useConversationStore.getState().setDraft("please clarify restore window")

    await useConversationStore.getState().submitComposerMessage()
    await waitForPollCycle()

    let state = useConversationStore.getState()
    expect(state.pendingTurnByConversationId.conv_clarifying_loop).toEqual({
      state: "clarifying",
      turnId: "turn_clarifying_001",
      sourceMessageId: "msg_clarification_001",
    })
    expect(state.messagesByConversationId.conv_clarifying_loop?.map((message) => message.messageId)).toEqual([
      "msg_user_clarifying_001",
      "msg_clarification_001",
    ])
    expect(listConversationEventsMock).toHaveBeenCalledTimes(1)

    await new Promise((resolve) => setTimeout(resolve, 1100))

    state = useConversationStore.getState()
    expect(listConversationEventsMock).toHaveBeenCalledTimes(1)
    expect(state.pendingTurnByConversationId.conv_clarifying_loop).toEqual({
      state: "clarifying",
      turnId: "turn_clarifying_001",
      sourceMessageId: "msg_clarification_001",
    })
  })

  it("drains paginated event pages before treating the turn as complete", async () => {
    createConversationMock.mockResolvedValue({
      conversation: {
        conversationId: "conv_paginated_events",
        title: "Paginated events",
        createdAt: "2026-04-23T09:00:00.000Z",
        updatedAt: "2026-04-23T09:00:00.000Z",
        interactionState: "thinking",
        activeTurnId: "turn_paginated_001",
      },
      message: {
        messageId: "msg_user_paginated_001",
        conversationId: "conv_paginated_events",
        turnId: "turn_paginated_001",
        role: "user",
        contentType: "text",
        content: "show the latest result",
        createdAt: "2026-04-23T09:00:00.000Z",
        status: "published",
      },
      statusEvent: messageCreatedStatusEvent(
        "conv_paginated_events",
        "turn_paginated_001",
        "msg_user_paginated_001",
      ),
      nextPollAfterMs: 0,
    })

    listConversationEventsMock
      .mockResolvedValueOnce({
        events: [
          {
            statusEventId: "evt_thought_paginated_001",
            conversationId: "conv_paginated_events",
            turnId: "turn_paginated_001",
            messageId: "msg_thought_paginated_001",
            eventType: "rich_content.created",
            sequence: 2,
            interactionState: "thinking",
            messageStatus: "responded",
            message: {
              messageId: "msg_thought_paginated_001",
              conversationId: "conv_paginated_events",
              turnId: "turn_paginated_001",
              parentMessageId: "msg_user_paginated_001",
              role: "assistant",
              contentType: "rich_content",
              content: "Collecting context",
              createdAt: "2026-04-23T09:00:05.000Z",
              status: "responded",
              richPayload: {
                kind: "thought",
                data: {
                  summary: "Collecting context",
                },
              },
            },
            createdAt: "2026-04-23T09:00:05.000Z",
          },
        ],
        nextCursor: "cursor_after_2",
        hasMore: true,
        latestSequence: 4,
        recommendedPollIntervalMs: 5000,
        interactionState: "completed",
      })
      .mockResolvedValueOnce({
        events: [
          {
            statusEventId: "evt_result_paginated_001",
            conversationId: "conv_paginated_events",
            turnId: "turn_paginated_001",
            messageId: "msg_result_paginated_001",
            eventType: "rich_content.created",
            sequence: 3,
            interactionState: "completed",
            messageStatus: "responded",
            message: {
              messageId: "msg_result_paginated_001",
              conversationId: "conv_paginated_events",
              turnId: "turn_paginated_001",
              parentMessageId: "msg_user_paginated_001",
              role: "assistant",
              contentType: "rich_content",
              content: "Generated final result mock-951",
              createdAt: "2026-04-23T09:00:06.000Z",
              status: "responded",
              richPayload: {
                kind: "result",
                data: {
                  title: "Mock result",
                  summary: "Generated final result mock-951",
                  resultId: "mock-951",
                  sourceMessage: "show the latest result",
                },
              },
            },
            createdAt: "2026-04-23T09:00:06.000Z",
          },
          {
            statusEventId: "evt_completed_paginated_001",
            conversationId: "conv_paginated_events",
            turnId: "turn_paginated_001",
            eventType: "interaction.status_changed",
            sequence: 4,
            interactionState: "completed",
            activeTurnId: null,
            completedTurnId: "turn_paginated_001",
            createdAt: "2026-04-23T09:00:07.000Z",
          },
        ],
        nextCursor: "cursor_after_4",
        hasMore: false,
        latestSequence: 4,
        recommendedPollIntervalMs: 5000,
        interactionState: "completed",
      })

    await useConversationStore.getState().hydrate()
    useConversationStore.getState().setDraft("show the latest result")

    await useConversationStore.getState().submitComposerMessage()
    await waitForPollCycle()
    await new Promise((resolve) => setTimeout(resolve, 25))

    const state = useConversationStore.getState()
    expect(listConversationEventsMock).toHaveBeenCalledTimes(2)
    expect(state.messagesByConversationId.conv_paginated_events?.map((message) => message.messageId)).toEqual([
      "msg_user_paginated_001",
      "msg_thought_paginated_001",
      "msg_result_paginated_001",
    ])
    expect(state.pendingTurnByConversationId.conv_paginated_events).toEqual({ state: "idle" })
    expect(state.eventCursorByConversationId.conv_paginated_events).toBe("cursor_after_4")
    expect(state.latestEventSequenceByConversationId.conv_paginated_events).toBe(4)
  })

  it("does not start event polling when opening a conversation that is already clarifying", async () => {
    getConversationDetailMock.mockResolvedValueOnce({
      conversationId: "conv_restore_clarifying",
      title: "订单数据库恢复",
      createdAt: "2026-04-23T09:00:00.000Z",
      updatedAt: "2026-04-23T09:00:20.000Z",
      interactionState: "clarifying",
      activeTurnId: "turn_restore_clarifying",
    })
    getConversationMessagesMock.mockResolvedValueOnce([
      {
        messageId: "msg_clarification_restore_001",
        conversationId: "conv_restore_clarifying",
        turnId: "turn_restore_clarifying",
        role: "assistant",
        contentType: "clarification",
        content: "请确认目标恢复时间窗口。",
        createdAt: "2026-04-23T09:00:20.000Z",
        status: "responded",
      },
    ])

    await useConversationStore.getState().selectConversation("conv_restore_clarifying")

    const state = useConversationStore.getState()
    expect(state.pendingTurnByConversationId.conv_restore_clarifying).toEqual({
      state: "clarifying",
      turnId: "turn_restore_clarifying",
    })
    expect(listConversationEventsMock).not.toHaveBeenCalled()
  })

  it("does not start event polling when message state snapshots already mark the turn completed", async () => {
    getConversationDetailMock.mockResolvedValueOnce({
      conversationId: "conv_restore_completed",
      title: "Restore completed",
      createdAt: "2026-04-23T09:00:00.000Z",
      updatedAt: "2026-04-23T09:00:20.000Z",
      interactionState: "executing",
      activeTurnId: "turn_restore_completed",
    })
    getConversationMessagesMock.mockResolvedValueOnce([
      {
        messageId: "msg_result_restore_001",
        conversationId: "conv_restore_completed",
        turnId: "turn_restore_completed",
        role: "assistant",
        contentType: "rich_content",
        content: "Generated three candidate recovery plans.",
        createdAt: "2026-04-23T09:00:20.000Z",
        status: "responded",
        richPayload: {
          kind: "layout_tree",
          data: {
            activity: {
              contract: "conversation.ui.layout-tree@1",
              blockId: "candidate_compare_001",
              ui: {
                id: "root",
                type: "stack",
              },
            },
            stateSnapshot: {
              interaction: {
                status: "completed",
              },
              selection: {
                required: true,
                selectedCandidateOptionId: "candidate_option_a",
                selectionLocked: false,
              },
            },
          },
        },
      },
    ])

    await useConversationStore.getState().selectConversation("conv_restore_completed")

    const state = useConversationStore.getState()
    expect(state.pendingTurnByConversationId.conv_restore_completed).toEqual({ state: "idle" })
    expect(state.detailsByConversationId.conv_restore_completed?.interactionState).toBe("completed")
    expect(state.detailsByConversationId.conv_restore_completed?.activeTurnId).toBeUndefined()
    expect(listConversationEventsMock).not.toHaveBeenCalled()
  })

  it("restores pending turn timing from the active turn user message when reopening an active conversation", async () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date("2026-04-28T08:40:00.000Z"))

    try {
      getConversationDetailMock.mockResolvedValueOnce({
        conversationId: "conv_restore_timer",
        title: "Restore timer",
        createdAt: "2026-04-28T08:20:00.000Z",
        updatedAt: "2026-04-28T08:29:47.192Z",
        interactionState: "executing",
        activeTurnId: "turn_restore_timer_001",
      })
      getConversationMessagesMock.mockResolvedValueOnce([
        {
          messageId: "msg_user_restore_timer_001",
          conversationId: "conv_restore_timer",
          turnId: "turn_restore_timer_001",
          role: "user",
          contentType: "text",
          content: "recommend a backup plan",
          createdAt: "2026-04-28T08:25:37.353Z",
          status: "responded",
        },
        {
          messageId: "msg_assistant_restore_timer_001",
          conversationId: "conv_restore_timer",
          turnId: "turn_restore_timer_001",
          role: "assistant",
          contentType: "rich_content",
          content: "working",
          createdAt: "2026-04-28T08:27:23.268Z",
          status: "responded",
          richPayload: {
            kind: "markdown",
            data: {
              markdown: "working",
            },
          },
        },
      ])

      await useConversationStore.getState().selectConversation("conv_restore_timer")

      const state = useConversationStore.getState()
      expect(state.pendingTurnByConversationId.conv_restore_timer).toEqual({
        state: "thinking",
        turnId: "turn_restore_timer_001",
      })
      expect(
        (state as unknown as { pendingTurnStartedAtMsByConversationId: Record<string, number> })
          .pendingTurnStartedAtMsByConversationId.conv_restore_timer,
      ).toBe(new Date("2026-04-28T08:25:37.353Z").getTime())
    } finally {
      useConversationStore.getState().activateLocalDraftWorkspace()
      vi.useRealTimers()
    }
  })

  it("stops polling the previously selected conversation after switching to a different conversation", async () => {
    getConversationDetailMock.mockImplementation(async (conversationId: string) => {
      if (conversationId === "conv_polling_old") {
        return {
          conversationId: "conv_polling_old",
          title: "Polling old",
          createdAt: "2026-04-23T09:00:00.000Z",
          updatedAt: "2026-04-23T09:00:20.000Z",
          interactionState: "executing" as const,
          activeTurnId: "turn_polling_old",
        }
      }

      if (conversationId === "conv_idle_new") {
        return {
          conversationId: "conv_idle_new",
          title: "Idle new",
          createdAt: "2026-04-23T09:05:00.000Z",
          updatedAt: "2026-04-23T09:05:20.000Z",
          interactionState: "completed" as const,
        }
      }

      throw new Error(`Unexpected conversationId ${conversationId}`)
    })
    getConversationMessagesMock.mockResolvedValue([])

    await useConversationStore.getState().selectConversation("conv_polling_old")

    expect(useConversationStore.getState().pendingTurnByConversationId.conv_polling_old).toEqual({
      state: "thinking",
      turnId: "turn_polling_old",
    })

    await useConversationStore.getState().selectConversation("conv_idle_new")
    await new Promise((resolve) => setTimeout(resolve, 1100))

    expect(listConversationEventsMock).not.toHaveBeenCalled()
    expect(useConversationStore.getState().selectedWorkspace).toEqual({
      kind: "conversation",
      conversationId: "conv_idle_new",
    })
  })

  it("marks the pending turn as error from interaction status events", async () => {
    sendMessageMock.mockResolvedValue({
      conversation: {
        conversationId: "conv_error_event",
        title: "Error event",
        createdAt: "2026-04-23T09:00:00.000Z",
        updatedAt: "2026-04-23T09:01:00.000Z",
        interactionState: "thinking",
        activeTurnId: "turn_error",
      },
      message: {
        messageId: "msg_user_error_001",
        conversationId: "conv_error_event",
        turnId: "turn_error",
        role: "user",
        contentType: "text",
        content: "check restore status",
        createdAt: "2026-04-23T09:01:00.000Z",
        status: "published",
      },
      statusEvent: messageCreatedStatusEvent("conv_error_event", "turn_error", "msg_user_error_001"),
      nextPollAfterMs: 0,
    })

    listConversationEventsMock.mockResolvedValueOnce({
      events: [
        {
          statusEventId: "evt_error_001",
          conversationId: "conv_error_event",
          turnId: "turn_error",
          eventType: "interaction.status_changed",
          sequence: 2,
          interactionState: "error",
          activeTurnId: null,
          createdAt: "2026-04-23T09:01:10.000Z",
        },
      ],
      nextCursor: "cursor_after_error",
      hasMore: false,
      latestSequence: 2,
      recommendedPollIntervalMs: 1000,
      interactionState: "error",
    })

    useConversationStore.setState({
      bootstrapped: true,
      selectedWorkspace: {
        kind: "conversation",
        conversationId: "conv_error_event",
      },
      detailsByConversationId: {
        conv_error_event: {
          conversationId: "conv_error_event",
          title: "Error event",
          createdAt: "2026-04-23T09:00:00.000Z",
          updatedAt: "2026-04-23T09:00:00.000Z",
          interactionState: "idle",
        },
      },
      messagesByConversationId: {
        conv_error_event: [],
      },
      draftsByKey: {
        [conversationDraftKeyForConversation("conv_error_event")]: "check restore status",
      },
    })

    await useConversationStore.getState().submitComposerMessage()
    await waitForPollCycle()

    const state = useConversationStore.getState()
    expect(state.pendingTurnByConversationId.conv_error_event).toEqual({
      state: "error",
      turnId: "turn_error",
    })
  })

  it("sends the second plain-text turn through the existing conversation id", async () => {
    sendMessageMock.mockResolvedValue({
      conversation: {
        conversationId: "conv_text_loop",
        title: "Text loop",
        createdAt: "2026-04-23T09:00:00.000Z",
        updatedAt: "2026-04-23T09:01:00.000Z",
        interactionState: "thinking",
        activeTurnId: "turn_text_002",
      },
      message: {
        messageId: "msg_user_002",
        conversationId: "conv_text_loop",
        turnId: "turn_text_002",
        role: "user",
        contentType: "text",
        content: "restore to yesterday afternoon",
        createdAt: "2026-04-23T09:01:00.000Z",
        status: "published",
      },
      statusEvent: messageCreatedStatusEvent("conv_text_loop", "turn_text_002", "msg_user_002"),
      nextPollAfterMs: 0,
    })

    useConversationStore.setState({
      bootstrapped: true,
      selectedWorkspace: {
        kind: "conversation",
        conversationId: "conv_text_loop",
      },
      detailsByConversationId: {
        conv_text_loop: {
          conversationId: "conv_text_loop",
          title: "Text loop",
          createdAt: "2026-04-23T09:00:00.000Z",
          updatedAt: "2026-04-23T09:00:20.000Z",
          interactionState: "idle",
        },
      },
      messagesByConversationId: {
        conv_text_loop: [
          {
            messageId: "msg_user_001",
            conversationId: "conv_text_loop",
            role: "user",
            contentType: "text",
            content: "help me restore the order database",
            createdAt: "2026-04-23T09:00:00.000Z",
            status: "persisted",
          },
          {
            messageId: "msg_assistant_001",
            conversationId: "conv_text_loop",
            role: "assistant",
            contentType: "text",
            content: "I will first inspect the latest restore points.",
            createdAt: "2026-04-23T09:00:20.000Z",
            status: "responded",
          },
        ],
      },
      draftsByKey: {
        [conversationDraftKeyForConversation("conv_text_loop")]: "restore to yesterday afternoon",
      },
    })

    await useConversationStore.getState().submitComposerMessage()

    expect(sendMessageMock).toHaveBeenCalledWith(
      "conv_text_loop",
      expect.objectContaining({
        type: "user_message",
        content: "restore to yesterday afternoon",
      }),
    )
  })
})
