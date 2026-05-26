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
    ArticleID: number;
    Username: string;
    Content: string;
    CreatedAt: string;
}