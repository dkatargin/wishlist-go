import * as React from 'react';
import {
    Container,
} from "@mui/material";
import {useEffect, useState} from "react";
import {useParams, useLocation, useNavigate} from "react-router";
import {CreateWishlist, EditWish} from "../api/api";
import Header from '../components/header';
import {CreateWishlistDialog} from "../components/dialogs";
import Wishlist from "../components/wishlist";

export default function WishlistPage() {
    const wishlistId = useParams().id;
    const location = useLocation();
    const is_favorites = location.pathname.includes('favorites');
    const [createListDialogOpen, setCreateListDialogOpen] = useState(false);
    const [hasData, setHasData] = useState(false);
    const navigate = useNavigate();
    const [editListData, setEditListData] = useState(undefined);

    useEffect(() => {
        if (!wishlistId) {
            if (location.pathname === "/wishlist/new") {
                setCreateListDialogOpen(true);
            } else {
                throw new Error("Wishlist not found")
            }
        } else {
            setHasData(true);
        }
    }, [location, wishlistId]);

    const click_callback = (id: string, click_type: string): any => {
        // Обработка клика по карточке желания
        if (click_type === "set_done") {
            return EditWish(wishlistId, id, {status: "purchased"})
        } else if (click_type === "set_undone") {
            return EditWish(wishlistId, id, {status: "pending"})
        } else if (click_type === "card") {
            return true;
        }
    }

    const createListHandler = () => {
        console.log("Create list", editListData);
        setCreateListDialogOpen(false)
        CreateWishlist(editListData?.name, editListData?.description || "").then((result) => {
            try {
                navigate(`/wishlist/${result.share_code}`);
            } catch (e) {
                console.error("Error creating wishlist", e);
            }
        })
    }

    const handleCloseDialog = () => {
        setCreateListDialogOpen(false);
        navigate(-1);
    }

    return (
        <>
            <Header title={"Желания"} addType={"wish"}/>
            <Container sx={{
                alignItems: 'center',
                justifyContent: 'center',
                paddingLeft: 0,
                paddingRight: 0,
                width: {
                    xs: '100%',
                    sm: '80%',
                    md: '70%',
                    lg: '70%',
                    xl: '70%'
                },
            }}>

                {hasData ?
                    <Wishlist is_favorites={is_favorites} wishlist_id={wishlistId}
                              click_callback={click_callback}/> : <></>}
                <CreateWishlistDialog dialogData={editListData} setDialogData={setEditListData}
                                      isOpen={createListDialogOpen} onClose={() => handleCloseDialog()}
                                      onConfirm={() => {
                                          createListHandler();
                                      }}/>
            </Container>
        </>
    )
}