interface List {
    id: number;
    owner: number;
    name: string;
    description: string;
    share_code: string;
    created_at: number;
    updated_at: number;
}

interface Wish {
    id: number;
    wishlist_code: string;
    owner_id: number;
    name: string;
    priority: number;
    status: string;
    market_link: string;
    market_picture: string;
    market_price: number;
    market_currency: string;
    market_quantity: number;
    created_at: number;
    updated_at: number;
}

interface Empty {}

export {List, Wish, Empty}