import {JSX, useEffect, useState} from "react";
import {useNavigate} from "react-router";
import {Card, CardActionArea, CardContent, Typography, Box} from "@mui/material";
import {Edit, DeleteForever} from "@mui/icons-material";
import {DeleteWishlist, EditWishlist, FetchFavorites, FetchLists} from "../api/api";
import {List} from "../api/interfaces"
import {CreateWishlistDialog} from "./dialogs";
import * as React from "react";


const WishlistCard = ({cardObject, onEdit, onDelete, onCardClick}: {
    cardObject: List,
    onEdit: any, onDelete: any, onCardClick: any
}) => {
    const cardActions = () => {
        return <>
            <Box

                onClick={(e) => {
                    e.stopPropagation();
                    onDelete(cardObject?.share_code);
                }}
                sx={{
                    position: 'absolute',
                    top: 8,
                    left: 8,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    width: '40px',
                    height: '40px',
                    borderRadius: '50%',
                    backgroundColor: '#f4f4f4',
                    cursor: 'pointer',
                    '&:hover': {
                        backgroundColor: '#d6d6d6',
                    },
                }}
            >
                <DeleteForever sx={{color: '#d02a2a', fontSize: '25px'}}/>
            </Box>
            <Box
                onClick={(e) => {
                    e.stopPropagation();
                    onEdit(cardObject);
                }}
                sx={{
                    position: 'absolute',
                    top: 8,
                    right: 8,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    width: '40px',
                    height: '40px',
                    borderRadius: '50%',
                    backgroundColor: '#f4f4f4',
                    cursor: 'pointer',
                    '&:hover': {
                        backgroundColor: '#d6d6d6',
                    },
                }}
            >
                <Edit sx={{color: '#424242', fontSize: '25px'}}/>
            </Box>
        </>
    }

    return <Card>
        <CardActionArea onClick={() => onCardClick(cardObject?.share_code)}

                        sx={{
                            height: '100%',
                            '&[data-active]': {
                                backgroundColor: 'action.selected',
                                '&:hover': {
                                    backgroundColor: 'action.selectedHover',
                                },
                            },
                        }}
        >
            {cardActions()}
            <CardContent sx={{
                display: 'flex',
                justifyContent: 'flex-end',
                alignItems: 'flex-end',
                height: '70%',

            }}>

                <Typography variant="h5" component="div">
                    {cardObject?.name}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    {cardObject?.description}
                </Typography>
            </CardContent>
        </CardActionArea>
    </Card>

}

// const ListCards = ({data, is_favorites, card_click_handler}:
//                    {
//                        data: List[],
//                        is_favorites: boolean,
//                        card_click_handler: (id: string) => void
//                    }): JSX.Element[] => {
//     let cards: JSX.Element[];
//
//     const handleCardEditList = (listData: object) => {
//
//         console.log(listData);
//         // Обработчик клика по карточке
//
//         // EditWishlist(id, [{name: "name", value: "new name"}]).then(() => {
//
//     }
//
//     const handleCardDeleteList = (id: string) => {
//         // Обработчик клика по карточке
//         DeleteWishlist(id).then(() => {
//             location.reload();
//         })
//     }
//
//     const handleCardClick = (id: string) => {
//         navigate(`/wishlist/${id}`);
//     }
//
//     const handleCardClick = (id: string) => {
//
//     }
//
//
//     // Отображает карточки списков
//     cards = Array.from(data).map((item: List) => (
//         WishlistCard(item, handleCardEditList, handleCardDeleteList, handleCardClick)
//         // @ts-ignore
//
//     ))
//     // Если is_favorites=true, то добавляем псевдокарточку для добавления нового списка
//     if (!is_favorites) {
//         cards.unshift(
//             <Card key={-1} sx={{height: "100%", display: 'flex', alignItems: 'center', justifyContent: 'center'}}>
//                 <CardActionArea
//                     onClick={() => card_click_handler("new")}
//                     sx={{height: "100%"}}>
//                     <CardContent>
//                         <Typography align={'center'} component="div" variant="h5">
//                             <Add sx={{width: '100px', height: '100px', color: '#e0e0e0'}}/>
//                         </Typography>
//                     </CardContent>
//                 </CardActionArea>
//             </Card>
//         )
//     }
//     // собираем все списки в массив
//     return cards;
// };

export default function CardsList({isEditable}: { isEditable: boolean }): JSX.Element {
    // Возвращает списки желаний для ListsPage
    // is_favorites - true, если нужно отобразить сохраненные чужие списки
    const [lists, setLists] = useState(Array<List>());
    const [dataLoaded, setDataLoaded] = useState(false);
    const [createListDialogOpen, setCreateListDialogOpen] = useState(false);
    const [editDialogData, setEditDialogData] = useState<List | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        if (dataLoaded) {
            return;
        }
        // Fetch lists and favorites from API
        const fetchFavorites = async () => {
            const data = await FetchFavorites();
            setLists(data);
        }

        const fetchLists = async () => {
            const data = await FetchLists();
            setLists(data);
        }

        if (isEditable) {
            fetchFavorites().then(() => {
                    setDataLoaded(true)
                }
            )
        } else {
            fetchLists().then(() => {
                    setDataLoaded(true)
                }
            )
        }
    }, [isEditable]);

    const cardClickHandler = (id: string) => {
        navigate(`/wishlist/${id}`);

    }

    const cardDeleteHandler = (id: string) => {
        const confirm = window.confirm("Удалить список?")
        if (!confirm) {
            return;
        }
        DeleteWishlist(id).then(() => {
            setLists(lists.filter((item: List) => item.share_code !== id));
        })
    }

    const dialogCloseHandler = () => {
        setCreateListDialogOpen(false);
        navigate(-1);
    }

    const cardEditHandler = (listData: List) => {
        setEditDialogData(listData);
        setCreateListDialogOpen(true);
    }

    const cardSaveHandler = () => {

        if (editDialogData) {
            // TODO: добавить проверку на пустое имя
            if (!editDialogData.name) {
                return;
            }
            EditWishlist(editDialogData.share_code, {
                name: editDialogData.name,
                description: editDialogData.description
            }).then(() => {
                setLists(lists.map((item: List) => {
                    if (item.share_code === editDialogData?.share_code) {
                        return editDialogData;
                    }
                    return item;
                }))
            })
        } else {
            console.error("No data to save");
        }

        setCreateListDialogOpen(false);
    }

    const cardsList = () => {
        return lists.map((item: List) => (
            <WishlistCard key={item.share_code} cardObject={item} onEdit={cardEditHandler} onDelete={cardDeleteHandler}
                          onCardClick={cardClickHandler}/>
        ))
    }

    return (
        <Box
            sx={{
                width: '100%',
                display: 'grid',
                gridTemplateColumns: 'repeat(2, minmax(min(100px, 100%), 1fr))',
                gridAutoRows: {
                    xs: 'minmax(150px, 150px )',
                    sm: 'minmax(250px, 250px )',
                    md: 'minmax(250px, 250px )',
                    lg: 'minmax(300px, 300px )',
                    xl: 'minmax(300px, 300px )'
                },
                gap: 2,
                columns: 2,
                marginTop: 2,
                alignContent: 'center',
                justifyContent: 'center',
            }}
        >
            {cardsList()}
            <CreateWishlistDialog dialogData={editDialogData} setDialogData={setEditDialogData}
                                  isOpen={createListDialogOpen} onClose={() => dialogCloseHandler()}
                                  onConfirm={() => {
                                      cardSaveHandler();
                                  }}/>
        </Box>
    );
}
