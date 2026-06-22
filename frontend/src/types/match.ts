export interface Team {
  id: string;
  name_en: string;
  name_fa: string;
  fifa_code: string;
  groups: string;
  flag: string;
}

export interface Scorer {
  player: string;
  minute: string;
  team_id: string;
}

export interface Match {
  _id: string;
  id: string;
  home_team_id: string;
  away_team_id: string;
  home_score: string;
  away_score: string;
  home_scorers: Scorer[] | null;
  away_scorers: Scorer[] | null;
  group: string;
  matchday: string;
  local_date: string;
  persian_date: string;
  stadium_id: string;
  finished: "TRUE" | "FALSE";
  time_elapsed: string;
  type: string;
  home_team_name_en?: string;
  home_team_name_fa?: string;
  away_team_name_en?: string;
  away_team_name_fa?: string;
  home_team_label?: string;
  away_team_label?: string;
}

export interface GoalEvent {
  player: string;
  minute: number;
  teamId: string;
}
