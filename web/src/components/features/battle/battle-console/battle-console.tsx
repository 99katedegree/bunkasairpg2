export type BattleConsoleProps = {
  children: React.ReactNode;
};

export function BattleConsole({ children }: BattleConsoleProps) {
  return (
    <div className="fixed bottom-0 w-full bg-black/70 text-white p-2">
      {children}
    </div>
  );
}
