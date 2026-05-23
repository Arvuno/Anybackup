import { createContext, useContext, useMemo, useState, type ReactNode } from "react"
import { getStoredLocale, LANGUAGE_STORAGE_KEY, messages, type Locale, type MessageKey } from "@/i18n/messages"

export const LANGUAGE_SWITCH_ENABLED = false

interface I18nContextValue {
  locale: Locale
  setLocale: (locale: Locale) => void
  t: (key: MessageKey) => string
}

const I18nContext = createContext<I18nContextValue | null>(null)

export function I18nProvider({ children }: { children: ReactNode }) {
  const [locale, setLocaleState] = useState<Locale>(() => getStoredLocale())

  const setLocale = (nextLocale: Locale) => {
    setLocaleState(nextLocale)
    window.localStorage.setItem(LANGUAGE_STORAGE_KEY, nextLocale)
  }

  const contextValue = useMemo<I18nContextValue>(() => {
    const bundle = messages[locale]
    return {
      locale,
      setLocale,
      t: (key) => bundle[key] ?? messages.en[key],
    }
  }, [locale])

  return <I18nContext.Provider value={contextValue}>{children}</I18nContext.Provider>
}

export function useI18n() {
  const context = useContext(I18nContext)
  if (!context) throw new Error("useI18n must be used within I18nProvider")
  return context
}

