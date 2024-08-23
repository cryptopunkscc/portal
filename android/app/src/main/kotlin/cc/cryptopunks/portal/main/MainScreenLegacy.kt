package cc.cryptopunks.portal.main

//import android.net.Uri
//import androidx.activity.compose.BackHandler
//import androidx.compose.material.icons.Icons
//import androidx.compose.material.icons.filled.Face
//import androidx.compose.material.icons.filled.List
//import androidx.compose.material.icons.filled.Person
//import androidx.compose.material.icons.filled.PlayArrow
//import androidx.compose.material.icons.filled.Settings
//import androidx.compose.runtime.Composable
//import androidx.compose.runtime.getValue
//import androidx.compose.runtime.mutableStateOf
//import androidx.compose.runtime.remember
//import androidx.compose.runtime.setValue
//import androidx.navigation.compose.NavHost
//import androidx.navigation.compose.composable
//import androidx.navigation.compose.rememberNavController
//import cc.cryptopunks.portal.admin.AdminScreen
//import cc.cryptopunks.portal.admin.AdminWrapToggle
//import cc.cryptopunks.portal.compose.AstralTheme
//import cc.cryptopunks.portal.config.ConfigEditorScreen
//import cc.cryptopunks.portal.config.ConfigScreen
//import cc.cryptopunks.portal.contacts.ContactsScreen
//import cc.cryptopunks.portal.dashboard.DashboardItem
//import cc.cryptopunks.portal.dashboard.DashboardScreen
//import cc.cryptopunks.portal.js.JsAppsScreen
//import cc.cryptopunks.portal.log.LogScreen
//import cc.cryptopunks.portal.log.LogWrapToggle
//
//@Composable
//fun MainScreen() {
//    AstralTheme {
//        val mainNavController = rememberNavController()
//        val dashboardNavController = rememberNavController()
//        var showBars by remember { mutableStateOf(true) }
//        NavHost(
//            navController = mainNavController, startDestination = "dashboard"
//        ) {
//            composable("dashboard") {
//                DashboardScreen(
//                    navController = dashboardNavController,
//                    actions = { MainServiceToggle() },
//                    showBars = showBars,
//                    items = remember {
//                        listOf(
//                            DashboardItem("log", Icons.Default.List,
//                                actions = { LogWrapToggle() }
//                            ) {
//                                LogScreen()
//                            },
//                            DashboardItem("config", Icons.Default.Settings) {
//                                ConfigScreen { file ->
//                                    val route = "config_editor/${Uri.encode(file.absolutePath)}"
//                                    mainNavController.navigate(route)
//                                }
//                            },
//                            DashboardItem("admin", Icons.Default.Face,
//                                actions = {
//                                    MainBarToggle { showBars = false }
//                                    AdminWrapToggle()
//                                }
//                            ) {
//                                AdminScreen()
//                            },
//                            DashboardItem("apps", Icons.Default.PlayArrow) {
//                                JsAppsScreen()
//                            },
//                            DashboardItem("contacts", Icons.Default.Person) {
//                                ContactsScreen()
//                            },
//                        )
//                    },
//                )
//                BackHandler(!showBars) {
//                    showBars = true
//                }
//            }
//            composable("config_editor/{file}") {
//                ConfigEditorScreen(mainNavController)
//            }
//        }
//    }
//}
