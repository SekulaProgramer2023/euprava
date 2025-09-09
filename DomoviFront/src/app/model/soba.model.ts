export interface Soba {
  id: string;
  roomNumber: string;
  capacity: number;
  users: string[];
  onBudget: boolean;
  IsFree: boolean; // koristi veliko I jer backend tako vraća
}
