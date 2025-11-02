import React from 'react';
import {useEffect} from 'react';
import {Container} from "@mui/material";
import {useLocation} from "react-router";
import Header from '../components/header';
import CardsList from "../components/lists";

export default function ListsPage() {
    const location = useLocation();
    const isFavorites = location.pathname.includes('favorites');

    useEffect(() => {
    }, [])


    return (
        <>
        <Header title={"Списки желаний"} addType={"list"}/>
        <Container sx={{
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            width: {
                xs: '100%',
                sm: '70%',
                md: '70%',
                lg: '70%',
                xl: '70%'
            }
        }}>

            {isFavorites ? <CardsList is_favorites={true}/> : <CardsList is_favorites={false}/>}
        </Container>
        </>
    )
}