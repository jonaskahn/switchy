export interface AlternateLaunch {
  Name: string
  LaunchArgs: string
}

export interface Browser {
  Name: string
  ExePath: string
  IconSpec: string
  LaunchArgs: string
  IsUwp: boolean
  Hidden: boolean
  AlternateLaunches?: AlternateLaunch[]
  ProfileIcon?: string
}

export interface Rule {
  Pattern: string
  Browser: string
}

export interface Ruleset {
  Name: string
  Enabled: boolean
  Rules: Rule[]
}

export interface AppSettings {
  UseRules: boolean
  TimedSelection: boolean
  TimedSeconds: number
  DefaultBrowser: string
  Appearance: string
  Theme: string
  SelectorLayout: string
}

export interface UserSettings {
  Browsers: Browser[]
  Rulesets: Ruleset[]
  AppSettings: AppSettings
}
