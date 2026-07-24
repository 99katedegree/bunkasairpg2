export type Boss = {
  id: string;
  name: string;
  imageUrl: string;
  attack: number;
  maxHitPoint: number;
  hitPoint: number;
  experiencePoint: number;
  slash: number;
  blow: number;
  shoot: number;
  neutral: number;
  flame: number;
  water: number;
  wood: number;
  shine: number;
  dark: number;
  weapon: null;
  item: null;
};
export const boss: Boss = {
  id: "990",
  name: "ヴァリアス",
  imageUrl: "/boss-neutral.png",
  attack: 1000,
  maxHitPoint: 4000,
  hitPoint: 4000,
  experiencePoint: 16000,
  slash: 1.0,
  blow: 1.0,
  shoot: 1.0,
  neutral: 1.0,
  flame: 1.4,
  water: 1.4,
  wood: 1.4,
  shine: 1.4,
  dark: 1.4,
  weapon: null,
  item: null,
};
