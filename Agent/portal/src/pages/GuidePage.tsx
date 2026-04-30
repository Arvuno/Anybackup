import { useCallback, useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import { AnimatePresence, motion } from "framer-motion"
import {
  ArrowRight,
  Bot,
  Check,
  ChevronRight,
  ClipboardCheck,
  Search,
  Shield,
  Sparkles,
  User,
  Zap,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { publicAssetUrlFromBasePath } from "@/config/base-path"
import { routes } from "@/config/routes"
import { useI18n } from "@/i18n"
import { cn } from "@/lib/cn"
import { dismissGuideForNow, markGuideDone, markGuideDontShowAgain } from "@/lib/session"

const logoIconUrl = publicAssetUrlFromBasePath(import.meta.env.BASE_URL, "/images/logo-icon.png")

const steps = ["welcome", "how-it-works", "demo-submit", "demo-plan", "demo-review", "ready"] as const

type GuideStep = (typeof steps)[number]

const workflowSteps = [
  { icon: Bot, labelKey: "guide.workflowStep1Label", subKey: "guide.workflowStep1Sub" },
  { icon: Search, labelKey: "guide.workflowStep2Label", subKey: "guide.workflowStep2Sub" },
  { icon: ClipboardCheck, labelKey: "guide.workflowStep3Label", subKey: "guide.workflowStep3Sub" },
  { icon: Zap, labelKey: "guide.workflowStep4Label", subKey: "guide.workflowStep4Sub" },
] as const

const productCards = [
  { name: "Anybackup Agent", descKey: "guide.product1Desc", color: "hsl(170 70% 45%)" },
  { name: "Anybackup Foundation", descKey: "guide.product2Desc", color: "hsl(205 90% 52%)" },
  { name: "Anybackup Client", descKey: "guide.product3Desc", color: "hsl(250 65% 58%)" },
] as const

const featureTags = [
  { textKey: "guide.feature1", color: "hsl(35 90% 55%)" },
  { textKey: "guide.feature2", color: "hsl(340 70% 55%)" },
  { textKey: "guide.feature3", color: "hsl(170 70% 45%)" },
  { textKey: "guide.feature4", color: "hsl(205 90% 52%)" },
  { textKey: "guide.feature5", color: "hsl(250 65% 58%)" },
] as const

export function GuidePage() {
  const { t } = useI18n()
  const navigate = useNavigate()
  const [step, setStep] = useState<GuideStep>("welcome")
  const [typeIndex, setTypeIndex] = useState(0)
  const [showAgent, setShowAgent] = useState(false)
  const [showPlan, setShowPlan] = useState(false)
  const [approved, setApproved] = useState(false)
  const [dontShowAgain, setDontShowAgain] = useState(false)

  const currentIndex = steps.indexOf(step)
  const userText = t("guide.userText")
  const agentText = t("guide.agentText")

  const goNext = useCallback(() => {
    const next = steps[currentIndex + 1]
    if (next) setStep(next)
  }, [currentIndex])

  useEffect(() => {
    if (step !== "demo-submit") return

    setTypeIndex(0)
    setShowAgent(false)
    const intervalId = window.setInterval(() => {
      setTypeIndex((value) => {
        if (value >= userText.length) {
          window.clearInterval(intervalId)
          window.setTimeout(() => setShowAgent(true), 600)
          return value
        }
        return value + 1
      })
    }, 45)

    return () => window.clearInterval(intervalId)
  }, [step])

  useEffect(() => {
    if (step !== "demo-submit" || !showAgent) return

    const timeoutId = window.setTimeout(() => setStep("demo-plan"), 2000)
    return () => window.clearTimeout(timeoutId)
  }, [showAgent, step])

  useEffect(() => {
    if (step !== "demo-plan") return

    setShowPlan(false)
    const timeoutId = window.setTimeout(() => setShowPlan(true), 800)
    return () => window.clearTimeout(timeoutId)
  }, [step])

  const enterWorkspace = () => {
    markGuideDone()
    navigate(routes.home, { replace: true })
  }

  const skipGuide = () => {
    if (dontShowAgain) {
      markGuideDontShowAgain()
    } else {
      dismissGuideForNow()
    }
    navigate(routes.home, { replace: true })
  }

  return (
    <div className="relative flex min-h-screen flex-1 flex-col items-center justify-center overflow-hidden bg-background">
      <div className="pointer-events-none absolute inset-0 bg-gradient-hero" />
      <div className="pointer-events-none absolute left-1/4 top-1/4 h-96 w-96 rounded-full bg-primary/[0.03] blur-3xl" />
      <div className="pointer-events-none absolute bottom-1/4 right-1/4 h-64 w-64 rounded-full bg-ai-purple/[0.03] blur-3xl" />

      <div className="absolute left-1/2 top-8 flex -translate-x-1/2 items-center gap-2" aria-label={t("guide.progressLabel")}>
        {steps.map((item, index) => (
          <div key={item} className="flex items-center gap-2">
            <div
              className={cn(
                "h-2 w-2 rounded-full transition-all duration-500",
                index <= currentIndex ? "scale-110 bg-primary" : "bg-border",
              )}
            />
            {index < steps.length - 1 && (
              <div
                className={cn(
                  "h-0.5 w-8 rounded-full transition-all duration-500",
                  index < currentIndex ? "bg-primary" : "bg-border",
                )}
              />
            )}
          </div>
        ))}
      </div>

      <AnimatePresence mode="wait">
        {step === "welcome" && (
          <motion.div
            key="welcome"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.5 }}
            className="w-full max-w-[1060px] px-6 text-center"
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              transition={{ type: "spring", stiffness: 200, damping: 15 }}
              className="mb-6"
            >
              <img src={logoIconUrl} alt="Anybackup" className="mx-auto mb-3 h-16 w-16 object-contain" />
              <h1 className="text-[28px] font-bold leading-none tracking-tight text-foreground">Anybackup</h1>
              <p className="mt-2 text-[9px] font-normal uppercase tracking-[0.25em] text-muted-foreground">
                Always Resilient
              </p>
            </motion.div>

            <h2 className="mb-2 text-lg font-bold text-foreground">{t("guide.howItWorksDescription")}</h2>

            <div className="mb-6">
              <div className="mx-auto mb-3 grid max-w-[1060px] gap-4 md:grid-cols-3">
                {productCards.map((product, index) => (
                  <motion.div
                    key={product.name}
                    initial={{ opacity: 0, y: 15 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.12 + 0.3 }}
                    className="rounded-xl border bg-card px-5 py-4 text-center shadow-card transition-all duration-300 hover:shadow-md"
                    style={{ borderColor: `${product.color}35` }}
                  >
                    <p className="mb-1.5 break-words text-sm font-bold text-primary">{product.name}</p>
                    <p className="break-words text-[11px] leading-relaxed text-foreground">{t(product.descKey)}</p>
                  </motion.div>
                ))}
              </div>

              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ delay: 0.7 }}
                className="mx-auto grid max-w-[1060px] gap-2 sm:grid-cols-5"
              >
                {featureTags.map((feature) => (
                  <span
                    key={feature.textKey}
                    className="whitespace-nowrap rounded-lg border py-1.5 text-center text-[11px] font-medium text-muted-foreground transition-all duration-200"
                    style={{ borderColor: `${feature.color}25`, backgroundColor: `${feature.color}06` }}
                  >
                    {t(feature.textKey)}
                  </span>
                ))}
              </motion.div>
            </div>

            <Button variant="ai" size="lg" onClick={goNext} className="gap-2 px-8">
              {t("guide.startLearn")} <ArrowRight className="h-4 w-4" />
            </Button>
          </motion.div>
        )}

        {step === "how-it-works" && (
          <motion.div
            key="how-it-works"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.5 }}
            className="max-w-2xl text-center"
          >
            <h2 className="mb-2 text-2xl font-bold text-foreground">{t("guide.howItWorksTitle")}</h2>
            <p className="mb-10 text-sm text-muted-foreground">{t("guide.howItWorksDescription")}</p>
            <div className="mb-10 flex items-start justify-center gap-4">
              {workflowSteps.map((workflow, index) => {
                const Icon = workflow.icon
                return (
                  <motion.div
                    key={workflow.labelKey}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.15 + 0.2 }}
                    className="relative flex w-36 flex-col items-center gap-3"
                  >
                    <div className="flex h-14 w-14 items-center justify-center rounded-xl border border-ai/10 bg-gradient-ai-subtle">
                      <Icon className="h-6 w-6 text-ai" />
                    </div>
                    <div>
                      <p className="text-sm font-semibold text-foreground">{t(workflow.labelKey)}</p>
                      <p className="mt-0.5 text-[11px] text-muted-foreground">{t(workflow.subKey)}</p>
                    </div>
                    {index < workflowSteps.length - 1 && (
                      <ChevronRight className="absolute left-[150px] top-[42px] h-4 w-4 text-border" />
                    )}
                  </motion.div>
                )
              })}
            </div>
            <Button variant="ai" size="lg" onClick={goNext} className="gap-2 px-8">
              {t("guide.viewDemo")} <ArrowRight className="h-4 w-4" />
            </Button>
          </motion.div>
        )}

        {step === "demo-submit" && (
          <motion.div
            key="demo-submit"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.5 }}
            className="w-full max-w-xl"
          >
            <h2 className="mb-1 text-center text-xl font-bold text-foreground">{t("guide.userDescribeTitle")}</h2>
            <p className="mb-6 text-center text-xs text-muted-foreground">{t("guide.userDescribeDescription")}</p>
            <div className="rounded-xl border border-border bg-card p-6 shadow-card">
              <div className="mb-4 flex flex-row-reverse gap-3">
                <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-secondary">
                  <User className="h-4 w-4 text-secondary-foreground" />
                </div>
                <div className="rounded-xl rounded-tr-sm bg-secondary px-4 py-3 text-sm text-secondary-foreground">
                  {userText.slice(0, typeIndex)}
                  {typeIndex < userText.length && (
                    <span className="ml-0.5 inline-block h-4 w-0.5 animate-blink align-middle bg-foreground" />
                  )}
                </div>
              </div>
              {showAgent && (
                <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} className="flex gap-3">
                  <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-gradient-ai shadow-ai">
                    <Bot className="h-4 w-4 text-primary-foreground" />
                  </div>
                  <div className="whitespace-pre-wrap rounded-xl rounded-tl-sm border border-border bg-card px-4 py-3 text-sm shadow-card">
                    {agentText}
                  </div>
                </motion.div>
              )}
            </div>
          </motion.div>
        )}

        {step === "demo-plan" && (
          <motion.div
            key="demo-plan"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.5 }}
            className="w-full max-w-xl"
          >
            <h2 className="mb-1 text-center text-xl font-bold text-foreground">{t("guide.agentGenerateTitle")}</h2>
            <p className="mb-6 text-center text-xs text-muted-foreground">{t("guide.agentGenerateDescription")}</p>
            {showPlan ? (
              <motion.div
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ type: "spring", stiffness: 300, damping: 25 }}
                className="overflow-hidden rounded-xl border border-border bg-card shadow-card"
              >
                <div className="p-5 pb-4">
                  <div className="mb-3 flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Shield className="h-5 w-5 text-ai" />
                      <h4 className="text-base font-semibold">{t("guide.planTitle")}</h4>
                    </div>
                    <span className="rounded-full bg-warning-surface px-2 py-0.5 text-[10px] font-medium text-warning">
                      {t("guide.pendingReview")}
                    </span>
                  </div>
                  <p className="mb-4 text-xs text-muted-foreground">{t("guide.planEnvironment")}</p>
                  <div className="grid grid-cols-2 gap-2">
                    {[
                      { label: t("guide.planMetricBackupType"), value: t("guide.planMetricBackupTypeValue") },
                      { label: t("guide.planMetricStrategy"), value: t("guide.planMetricStrategyValue") },
                      { label: t("guide.planMetricRpoRto"), value: t("guide.planMetricRpoRtoValue") },
                      { label: t("guide.planMetricDataReduction"), value: t("guide.planMetricDataReductionValue") },
                    ].map((item) => (
                      <div key={item.label} className="rounded-lg bg-secondary/60 px-3 py-2.5">
                        <p className="text-[10px] text-muted-foreground">{item.label}</p>
                        <p className="mt-0.5 text-xs font-medium text-foreground">{item.value}</p>
                      </div>
                    ))}
                  </div>
                </div>
                <div className="px-5 pb-4">
                  <div className="flex items-start gap-2 rounded-lg border border-ai/10 bg-ai-surface/70 p-2.5">
                    <Sparkles className="mt-0.5 h-3.5 w-3.5 shrink-0 text-ai" />
                    <p className="text-[11px] text-accent-foreground">{t("guide.planRecommendation")}</p>
                  </div>
                </div>
                <div className="border-t border-border p-4 text-center">
                  <Button variant="ai" onClick={goNext} className="gap-2 px-6">
                    {t("guide.nextReview")} <ArrowRight className="h-4 w-4" />
                  </Button>
                </div>
              </motion.div>
            ) : (
              <div className="flex flex-col items-center py-16">
                <div className="mb-4 flex h-12 w-12 animate-pulse-ai items-center justify-center rounded-xl bg-gradient-ai shadow-ai-lg">
                  <Bot className="h-6 w-6 text-primary-foreground" />
                </div>
                <p className="text-sm text-muted-foreground">{t("guide.agentGeneratingPlan")}</p>
              </div>
            )}
          </motion.div>
        )}

        {step === "demo-review" && (
          <motion.div
            key="demo-review"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.5 }}
            className="w-full max-w-md text-center"
          >
            <h2 className="mb-1 text-xl font-bold text-foreground">{t("guide.reviewTitle")}</h2>
            <p className="mb-8 text-xs text-muted-foreground">{t("guide.reviewDescription")}</p>
            {!approved ? (
              <motion.div
                initial={{ scale: 0.95 }}
                animate={{ scale: 1 }}
                className="mb-4 rounded-xl border border-border bg-card p-6 shadow-card"
              >
                <div className="mb-4 flex items-center gap-2">
                  <Shield className="h-5 w-5 text-ai" />
                  <span className="font-semibold">{t("guide.planTitle")}</span>
                </div>
                <p className="mb-6 text-xs text-muted-foreground">{t("guide.reviewPlanSummary")}</p>
                <Button variant="ai" size="lg" className="w-full gap-2" onClick={() => setApproved(true)}>
                  <Check className="h-5 w-5" />
                  {t("guide.reviewPassed")}
                </Button>
              </motion.div>
            ) : (
              <motion.div
                initial={{ opacity: 0, scale: 0.9 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ type: "spring" }}
                className="mb-4"
              >
                <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-success-surface">
                  <Check className="h-8 w-8 text-success" />
                </div>
                <p className="mb-2 text-lg font-semibold text-success">{t("guide.planApproved")}</p>
                <p className="mb-8 text-sm text-muted-foreground">{t("guide.agentExecuteBackup")}</p>
                <Button variant="ai" size="lg" onClick={goNext} className="gap-2 px-8">
                  {t("guide.finishGuide")} <ArrowRight className="h-4 w-4" />
                </Button>
              </motion.div>
            )}
          </motion.div>
        )}

        {step === "ready" && (
          <motion.div
            key="ready"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.5 }}
            className="text-center"
          >
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ type: "spring", stiffness: 200, delay: 0.2 }}
              className="mx-auto mb-8 flex h-20 w-20 items-center justify-center rounded-xl bg-gradient-ai shadow-ai-lg"
            >
              <Zap className="h-10 w-10 text-primary-foreground" />
            </motion.div>
            <h1 className="mb-2 text-2xl font-bold text-foreground">{t("guide.readyTitle")}</h1>
            <p className="mx-auto mb-10 whitespace-nowrap text-sm leading-relaxed text-muted-foreground">
              {t("guide.readyDescription")}
            </p>
            <Button variant="ai" size="lg" onClick={enterWorkspace} className="gap-2 px-10 text-base">
              {t("guide.enterWorkspace")} <ArrowRight className="h-5 w-5" />
            </Button>
          </motion.div>
        )}
      </AnimatePresence>

      {step !== "ready" && (
        <div className="absolute bottom-8 flex flex-col items-center gap-3">
          <label className="flex items-center gap-2 text-xs text-muted-foreground">
            <input
              type="checkbox"
              className="h-3.5 w-3.5 rounded border-input accent-[hsl(var(--ai))]"
              checked={dontShowAgain}
              onChange={(event) => setDontShowAgain(event.target.checked)}
            />
            {t("guide.dontShowAgain")}
          </label>
          <button
            type="button"
            onClick={skipGuide}
            className="text-xs text-muted-foreground transition-fast hover:text-foreground"
          >
            {t("guide.skipGuide")} →
          </button>
        </div>
      )}
    </div>
  )
}
