export interface Article {
    ID: string;
    Title: string;
    Preview: string;
    Content: string;
}

export interface Like{
    likes: number
}

export interface Comment {
    ID: number;
    CreatedAt: string;
    article_id: number;
    username: string;
    content: string;
}