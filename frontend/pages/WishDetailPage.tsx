import Header from "../components/header";
import React, {useEffect, useState} from "react";
import {
    Box,
    Button,
    Container,
    Stack,
    TextField,
    Alert,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions
} from "@mui/material";
import {useLocation, useParams, useNavigate} from "react-router";
import {CreateWish, EditWish, FetchWish, DeleteWish} from "../api/api";
import {Error} from "@mui/icons-material"


function DeleteConfirmDialog({open, onClose, onConfirm}) {
    return (
        <Dialog open={open} onClose={onClose}>
            <DialogTitle>Подтвердите удаление</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    Вы уверены, что хотите удалить это желание?
                </DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Отмена</Button>
                <Button onClick={onConfirm} color="error">Удалить</Button>
            </DialogActions>
        </Dialog>
    );
}

const WishDetailForm = ({data, onSubmit, onDelete}: {
    data: any | null,
    onDelete: any | null
    onSubmit: (event: React.FormEvent<HTMLFormElement>) => void,
}) => {
    const [productUrl, setProductUrl] = useState("");
    const [title, setTitle] = useState("");
    const [pictureUrl, setPictureUrl] = useState("");
    const [productPrice, setProductPrice] = useState("");
    const navigate = useNavigate();

    useEffect(() => {
        if (data) {
            setProductUrl(data.market_link);
            setTitle(data.name);
            setPictureUrl(data.market_picture);
            setProductPrice(data.market_price);
        }
    }, [data])

    return <Box component={"form"} sx={{marginTop: 2}} onSubmit={(e) => {
        onSubmit(e)
    }}>
        <Stack spacing={2}>

            <Stack direction="row" spacing={1} alignItems="center">
                <TextField
                    label="Ссылка на товар"
                    name="product_url"
                    required
                    fullWidth
                    value={productUrl}
                    onChange={(e) => {
                        setProductUrl(e.target.value)
                    }}
                />
                {!data && (<Button
                    variant="outlined"
                    type="submit"
                    formNoValidate
                    sx={{ whiteSpace: 'nowrap', px: '2rem', height: '3.5rem' }}
                >
                    Загрузить (beta)
                </Button>)}
            </Stack>

            <TextField
                label="Название желания"
                name="title"
                variant="outlined"
                fullWidth
                required
                value={title}
                onChange={(e) => {
                    setTitle(e.target.value)
                }}
            />
            <TextField
                label="URL картинки"
                name="picture_url"
                variant="outlined"
                fullWidth
                value={pictureUrl}
                onChange={(e) => {
                    setPictureUrl(e.target.value)
                }}
            />

            <TextField
                type="number"
                label="Цена"
                name="product_price"
                variant="outlined"
                fullWidth
                required
                value={productPrice}
                onChange={(e) => {
                    setProductPrice(e.target.value)
                }}
            />
            <Stack direction="row" spacing={1} sx={{mt: 1}}>
                <Button type="submit" variant="contained" fullWidth>Сохранить</Button>
                <Button variant="outlined" onClick={() => {
                    navigate(-1)
                }} fullWidth>Отмена</Button>
            </Stack>
            {data && (
                <Stack direction="row" spacing={1} sx={{mt: 1}}>
                    <Button variant="outlined" color="error" fullWidth
                            onClick={() => onDelete(data.id, false)}>Удалить</Button>
                </Stack>)}

        </Stack>
    </Box>
}

export default function WishDetailsPage() {
    const params = useParams();
    if (!params.id) {
        return
    }

    const location = useLocation();
    const navigate = useNavigate();

    const [wish, setWish] = useState<any | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [pageType, setPageType] = useState("");
    const [headerLabel, setHeaderLabel] = useState("");

    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);


    useEffect(() => {
        if (location.pathname.includes("wish/new")) {
            setHeaderLabel("Создание желания")
            setPageType("new")
            setIsLoading(false)
        } else {
            setHeaderLabel("Редактирование желания")
            setPageType("edit")
        }
    }, [location.pathname]);

    useEffect(() => {
        if (pageType === "edit" && params.id && params.wishid) {
            FetchWish(params.id, params.wishid).then((response) => {
                setWish(response);
                setIsLoading(false);
            });
        }
    }, [params.id, params.wishid, pageType]);

    const createWishFormHandler = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const formData = new FormData(event.currentTarget);
        const productPrice = Number(parseFloat(formData.get("product_price") as string).toFixed(2))

            CreateWish(params.id,
            formData.get("title"),
            formData.get("product_url"),
            formData.get("picture_url"),
            productPrice,
            "RUB", 1, 0
        ).then(() => {
            navigate(`/wishlist/${params.id}`);
        })
    }

    const editWishFormHandler = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const formData = new FormData(event.currentTarget);
        const productPrice = Number(parseFloat(formData.get("product_price") as string).toFixed(2))

        EditWish(params.id, params.wishid,
            {
                "name": formData.get("title"),
                "market_link": formData.get("product_url"),
                "market_picture": formData.get("picture_url"),
                "market_price": productPrice,
            }
        ).then(() => {
            navigate(`/wishlist/${params.id}`);
        })
    }

    const deleteWishFormHandler = (wishId: string, isConfirmed: boolean) => {
        if (!isConfirmed) {
            setDeleteDialogOpen(true);
        } else {
            DeleteWish(params.id, wishId).then(() => {
                navigate(`/wishlist/${params.id}`);
            })
        }
    }

    return (
        <>
            <Header title={headerLabel} addType={null}/>
            <Container sx={{
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                width: {
                    xs: '100%',
                    sm: '70%',
                    md: '70%',
                    lg: 700,
                    xl: 700
                },
                marginTop: 2,
                marginBottom: 2
            }}>

                {pageType === "new" ? (<WishDetailForm onSubmit={createWishFormHandler} data={null} onDelete={null}/>) :
                    (!isLoading ?
                        <WishDetailForm onDelete={deleteWishFormHandler} onSubmit={editWishFormHandler} data={wish}/> :
                        <Alert icon={<Error fontSize="inherit"/>} severity={"error"}>Не удалось загрузить
                            данные</Alert>)}
                <DeleteConfirmDialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}
                                     onConfirm={() => {
                                         deleteWishFormHandler(wish.id, true)
                                     }}/>
            </Container>
        </>
    );
}