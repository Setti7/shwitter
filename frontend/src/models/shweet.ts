import User from "./user";

export type Timeline = Shweet[];

export default interface Shweet {
    id: string;
    message: string;
    user: User;
    created_at: Date;
}