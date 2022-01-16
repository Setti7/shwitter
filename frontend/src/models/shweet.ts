import User from "./user";

export type Timeline = Shweet[];

export default interface Shweet {
  id: string;
  message: string;
  user: User;
  created_at: Date;
}

export interface ShweetDetails extends Shweet {
  like_count: number;
  reshweet_count: number;
  comment_count: number;
  liked: boolean;
  reshweeted: boolean;
}
