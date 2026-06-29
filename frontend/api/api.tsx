import {retrieveRawInitData} from '@telegram-apps/sdk';
import {List, Wish} from './interfaces';


const backendAPI = async (method: string, endpoint: string, body: object | null) => {
    let initDataRaw;
    try {
        initDataRaw = retrieveRawInitData()
    } catch (e) {
        if (import.meta.env.VITE_DEPLOYMENT_TYPE === 'development') {
            initDataRaw = import.meta.env.VITE_AUTH_MOCKUP
        } else {
            console.error("Error retrieving init data:", e);
            throw new Error("Error retrieving init data");
        }
    }

    console.log(import.meta.env.VITE_BACKEND_HOST)
    console.log(import.meta.env.VITE_BACKEND_PORT)
    console.log(import.meta.env.VITE_BACKEND_SCHEME)
    console.log(import.meta.env.VITE_AUTH_MOCKUP)
    console.log(import.meta.env.VITE_DEPLOYMENT_TYPE)

    const host = import.meta.env.VITE_BACKEND_HOST
    const port = import.meta.env.VITE_BACKEND_PORT
    const hostScheme = import.meta.env.VITE_BACKEND_SCHEME
    const url = `${hostScheme}://${host}:${port}/api/v1/${endpoint}`


    const headers = new Headers({
        'Content-Type': 'application/json',
        'Authorization': `tma ${initDataRaw}`
    });

    const response = await fetch(url, {
        method: method.toUpperCase(),
        headers: headers,
        body: method !== 'GET' ? JSON.stringify(body) : null
    }).catch((error) => {
        console.error("Error in fetch:", error);
    })

    if (!response) {
        throw new Error("No response from server");
    }

    if (method != "DELETE") {
        return response.json();
    }

};

const AddToFavorites = async (id: string): Promise<List[]> => {
    return await backendAPI("POST", `wishlist/${id}/favorite`, null);
}

// Получение всех списков
const FetchLists = async (): Promise<List[]> => {
    const response = await backendAPI("GET", "list", null);
    return response.wishlists || [];
};

// Получение "чужих" списков (избранное)
const FetchFavorites = async (): Promise<List[]> => {
    return await backendAPI("GET", "favorites", null);
}

// Создание нового списка
const CreateWishlist = async (name: string, description: string): Promise<List> => {
    const response = await backendAPI("POST", "list", {name: name, description: description});
    return response.wishlist;
}

// Получение списка желаний по ID
const FetchWishlist = async (id: string): Promise<Wish[]> => {
    const response = await backendAPI("GET", `list/${id}/wishes`, null);
    return response.wish_items || [];
}

// Редактирование списка
const EditWishlist = async (id: string, fields: object): Promise<List> => {
    const response = await backendAPI("PATCH", `list/${id}`, fields);
    return response.wishlist;
}

// Удаление списка
const DeleteWishlist = async (id: string): Promise<void> => {
    await backendAPI("DELETE", `list/${id}`, null);
}

const CreateWishWithCrawler = async (wishlistId: string, url: string): Promise<Wish> => {
    const response = await backendAPI("POST", `list/${wishlistId}/wishes/crawl`, {url: url});
    return response.wish_item;
}

// Создание желания в списке
const CreateWish = async (wishlistId: string,
                          name: string, market_url: string, market_pic: string, market_price: number,
                          market_currency: string, market_quantity: number, priority: number): Promise<Wish> => {

    if (market_url && !name) {
        return await CreateWishWithCrawler(wishlistId, market_url);
    }

    const response = await backendAPI("POST", `list/${wishlistId}/wishes`,
        {
            name: name, market_link: market_url, market_picture: market_pic, market_price: market_price,
            market_currency: market_currency, market_quantity: market_quantity, priority: priority
        });
    return response.wish_item;
}

// Получение конкретного желания из списка
const FetchWish = async (wishlistId: string, wishId: string):Promise<Wish> => {
    const response = await backendAPI("GET", `list/${wishlistId}/wishes/${wishId}`, null);
    return response.wish_item;
}

// Редактирование желания в списке
const EditWish = async (wishlistId: string, wishId: string, fields: object):Promise<Wish> => {
    const response = await backendAPI("PATCH", `list/${wishlistId}/wishes/${wishId}`, fields);
    return response.wish_item;
}

// Удаление желания из списка
const DeleteWish = async (wishlistId: string, wishId: string) :Promise<void> => {
    await backendAPI("DELETE", `list/${wishlistId}/wishes/${wishId}`, null);
}

export {AddToFavorites, FetchLists, FetchFavorites, CreateWishlist, FetchWishlist, EditWishlist,
    DeleteWishlist, CreateWish, FetchWish, EditWish, DeleteWish};