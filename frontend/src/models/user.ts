export default interface User {
  id: string;
  username: string;
  name: string;
  email: string;
  bio: string;
  joined_at: Date;
}

export interface UserProfile extends User {
  followers_count: number;
  friends_count: number;
  shweets_count: number;
}
