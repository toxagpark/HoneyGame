package menu_tg_transport

// callback data меню: menu:<action>, ставки — menu:bet:<percent>
const (
	callbackMenu         = "menu:"
	callbackMenuProfile  = "menu:profile"
	callbackMenuFight    = "menu:fight"
	callbackMenuFights   = "menu:fights"
	callbackMenuMyCall   = "menu:my-challenge"
	callbackMenuBack     = "menu:back"
	callbackMenuBetFmt   = "menu:bet:"      // + процент
	fightAcceptPrefix    = "fight:accept:"  // кнопки взятия вызова в меню
)
