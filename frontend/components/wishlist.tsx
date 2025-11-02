import {JSX, useEffect, useState} from "react";
import {useNavigate} from "react-router";
import {Box, Button, Card, CardActionArea, CardContent, Typography} from "@mui/material";
import {FetchWishlist} from "../api/api";
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import BookmarkAddIcon from '@mui/icons-material/BookmarkAdd';
import {Wish} from "../api/interfaces";


export default function Wishlist({wishlist_id, is_favorites, click_callback}:
                                 {
                                     wishlist_id: string,
                                     is_favorites: boolean,
                                     click_callback: (id: string, click_type: string) => any
                                 }): JSX.Element {
    // Возвращает список желаний для WishlistPage
    // is_personal - true, нужно для корректного отображения контролов
    const [wlData, setWlData] = useState(Array<Wish>());
    const navigate = useNavigate();

    const clickHandler = (id: string, click_type: string): any => {
        // Обработка клика по карточке желания
        if (click_type === "set_done") {
            click_callback(id, "set_done").then((data: any) => {
                    // @ts-ignore
                    setWlData(wlData.map(wish => wish.id.toString() === id ? data : wish));
                }
            )
        } else if (click_type === "set_undone") {
            click_callback(id, "set_done").then((data: any) => {
                    // @ts-ignore
                    setWlData(wlData.map(wish => wish.id.toString() === id ? data : wish));
                }
            )
        } else if (click_type === "card") {
            navigate(`/wishlist/${wishlist_id}/wish/${id}`);
        }
    }

    const cardChangeStateButton = (wish: Wish) => {
        if (is_favorites) {
            // если is_favorites то отобразить button_bookmark иначе отобразить кнопку done
            return <Button variant="contained" color="primary" size="small" sx={{marginRight: 1}}
                           onClick={() => clickHandler(wish.id.toString(), "set_hold")}>
                <BookmarkAddIcon/>
            </Button>
        } else {
            if (wish.status === "purchased") {
                return <Button variant="outlined" color="success" size="small" sx={{marginRight: 1}}
                               onClick={() => clickHandler(wish.id.toString(), "set_undone")}>
                    <CheckCircleIcon/>
                </Button>
            } else {
                return <Button variant="contained" color="primary" size="small"
                               sx={{
                                   marginRight: 1,
                                   alignItems: 'center',
                                   justifyContent: 'center'
                               }}
                               onClick={() => clickHandler(wish.id.toString(), "set_done")}>
                    <CheckCircleIcon/>
                </Button>
            }
        }
    }

    const card = (wish: Wish) => (
        <Card
            key={`wish-${wish.id}`}
            sx={{
                margin: 2,
                height: 150,
                display: 'flex',
                flexDirection: 'row',
                alignItems: 'center',
                borderRadius: 4, // Скругленные углы
                boxShadow: '0 4px 10px rgba(0, 0, 0, 0.1)', // Тень
                transition: 'transform 0.2s, box-shadow 0.2s', // Анимация
                '&:hover': {
                    transform: 'scale(1.02)', // Увеличение при наведении
                    boxShadow: '0 6px 15px rgba(0, 0, 0, 0.2)', // Более сильная тень
                },
            }}
        >


            <CardActionArea onClick={() => clickHandler(wish.id.toString(), "card")} sx={{
                width: '35%',
                height: '100%',
                textDecoration: 'none',
            }}>
                <Box
                    component="img"
                    src={wish.market_picture}
                    alt={wish.name}
                    sx={{
                        width: 100,
                        height: 100,
                        objectFit: 'cover',
                        margin: 2,
                    }}
                />
            </CardActionArea>
            <CardContent
                sx={{
                    display: 'flex',
                    flexDirection: 'column',
                    justifyContent: 'space-between',
                    flexGrow: 1,
                    padding: 2,
                    height: '70%',
                }}
            >
                <CardActionArea onClick={() => clickHandler(wish.id.toString(), "card")} sx={{
                    display: 'flex',
                    flexDirection: 'row',
                    alignItems: 'center',
                    width: '100%',
                    height: '100%',
                    textDecoration: 'none',
                }}>
                    <Typography variant="h6" component="div" sx={{
                        marginBottom: '1rem',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        fontSize: '20px',
                        lineHeight: '1.2',
                    }}>
                        {wish.name.length > 25 ? `${wish.name.slice(0, 20)}...` : wish.name}
                    </Typography>
                </CardActionArea>
                <Box
                    sx={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'flex-center',
                        marginTop: 'auto',
                    }}
                >
                    <Typography variant="body2" sx={{fontWeight: 'bold'}}>
                        {wish.market_price.toLocaleString('ru-RU')} {wish.market_currency}
                    </Typography>
                    {cardChangeStateButton(wish)}

                </Box>
            </CardContent>

        </Card>
    );


    useEffect(() => {
        const fetchData = async () => {
            const data = await FetchWishlist(wishlist_id);
            setWlData(data);
        }
        if (!wlData.length) {
            fetchData()
        }

    }, [wishlist_id, is_favorites])

    return (
        <Box>
            {wlData.map((wish) => (
                card(wish)
            ))}
        </Box>
    )


}