import React from "react";
import { Match, Team, GoalEvent, Scorer } from "../types/match";
import { mockTeams } from "../data/mockMatches";

interface MatchFeaturedCardProps {
  match: Match;
}

const MatchFeaturedCard: React.FC<MatchFeaturedCardProps> = ({ match }) => {
  const homeTeam = mockTeams.find((team) => team.id === match.home_team_id);
  const awayTeam = mockTeams.find((team) => team.id === match.away_team_id);

  if (!homeTeam || !awayTeam) {
    return <div className="text-red-500">Erro: Times não encontrados.</div>;
  }

  const getGoalEvents = (
    scorers: Scorer[] | null,
    teamId: string,
  ): GoalEvent[] => {
    if (!scorers) return [];
    return scorers.map((scorer) => ({
      player: scorer.player,
      minute: parseInt(scorer.minute),
      teamId: scorer.team_id,
    }));
  };

  const allGoalEvents: GoalEvent[] = [
    ...getGoalEvents(match.home_scorers, homeTeam.id),
    ...getGoalEvents(match.away_scorers, awayTeam.id),
  ].sort((a, b) => a.minute - b.minute);

  const firstHalfGoals = allGoalEvents.filter((event) => event.minute <= 45);
  const secondHalfGoals = allGoalEvents.filter((event) => event.minute > 45);

  const renderGoal = (goal: GoalEvent) => (
    <div
      key={`${goal.player}-${goal.minute}`}
      className="flex items-center space-x-1"
    >
      <span>⚽</span>
      <span>
        {goal.player} ({goal.minute}')
      </span>
    </div>
  );

  return (
    <div className="bg-[#1a1d2e] rounded-lg shadow-lg p-6 text-white max-w-md mx-auto my-8">
      {match.finished === "TRUE" && (
        <div className="flex justify-center mb-4">
          <span className="bg-emerald-600 px-3 py-1 text-xs font-bold rounded text-white">
            ENCERRADA
          </span>
        </div>
      )}

      <div className="flex justify-between items-center mb-6">
        <div className="flex flex-col items-center w-1/3">
          <img
            src={homeTeam.flag}
            alt={homeTeam.name_en}
            className="w-12 h-12 mb-2"
          />
          <span className="text-lg font-semibold text-center">
            {homeTeam.name_en}
          </span>
        </div>
        <div className="flex items-center justify-center w-1/3">
          <span className="font-bold text-5xl">{match.home_score}</span>
          <span className="font-bold text-5xl mx-4">-</span>
          <span className="font-bold text-5xl">{match.away_score}</span>
        </div>
        <div className="flex flex-col items-center w-1/3">
          <img
            src={awayTeam.flag}
            alt={awayTeam.name_en}
            className="w-12 h-12 mb-2"
          />
          <span className="text-lg font-semibold text-center">
            {awayTeam.name_en}
          </span>
        </div>
      </div>

      {(firstHalfGoals.length > 0 || secondHalfGoals.length > 0) && (
        <div className="mt-4 pt-4 border-t border-gray-700">
          {firstHalfGoals.length > 0 && (
            <div className="mb-4">
              <p className="text-center text-gray-400 text-sm mb-2">1T</p>
              <div className="flex flex-col space-y-1">
                {firstHalfGoals.map((goal) => (
                  <div
                    key={`${goal.player}-${goal.minute}`}
                    className={`flex ${goal.teamId === homeTeam.id ? "justify-start" : "justify-end"}`}
                  >
                    {renderGoal(goal)}
                  </div>
                ))}
              </div>
            </div>
          )}

          {secondHalfGoals.length > 0 && (
            <div>
              <p className="text-center text-gray-400 text-sm mb-2">2T</p>
              <div className="flex flex-col space-y-1">
                {secondHalfGoals.map((goal) => (
                  <div
                    key={`${goal.player}-${goal.minute}`}
                    className={`flex ${goal.teamId === homeTeam.id ? "justify-start" : "justify-end"}`}
                  >
                    {renderGoal(goal)}
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default MatchFeaturedCard;
