package cc.cryptopunks.portal.dashboard

import androidx.compose.foundation.layout.RowScope
import androidx.compose.foundation.layout.padding
import androidx.compose.material.BottomNavigation
import androidx.compose.material.BottomNavigationItem
import androidx.compose.material.Icon
import androidx.compose.material.Scaffold
import androidx.compose.material.Text
import androidx.compose.material.TopAppBar
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.PlayArrow
import androidx.compose.material.icons.filled.Settings
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableIntStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.tooling.preview.Preview
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController

data class DashboardItem(
    override val label: String,
    override val icon: ImageVector,
    override val actions: @Composable RowScope.() -> Unit = {},
    override val content: @Composable () -> Unit,
) : DashboardBottomItem, DashboardScreenItem {

    val route: String = label
}

interface DashboardBottomItem {
    val label: String
    val icon: ImageVector
}

interface DashboardScreenItem {
    val actions: @Composable RowScope.() -> Unit
    val content: @Composable () -> Unit
}


@Preview
@Composable
fun DashboardScreenPreview() {
    DashboardScreen(
        items = remember {
            listOf(
                DashboardItem("config", Icons.Default.Settings) {},
                DashboardItem("apps", Icons.Default.PlayArrow) {},
            )
        },
        showBars = true
    )
}


@Composable
fun DashboardScreen(
    navController: NavHostController = rememberNavController(),
    items: List<DashboardItem>,
    actions: @Composable (RowScope.() -> Unit) = { },
    showBars: Boolean,
) {
    var selected by rememberSaveable { mutableIntStateOf(0) }
    val onSelect = { item: DashboardItem ->
        val currentRoute = navController.currentDestination?.route
        if (currentRoute != null)
            navController.navigate(item.route) {
                popUpTo(currentRoute) {
                    inclusive = true
                }
            }
    }
    Scaffold(
        topBar = {
            if (showBars) TopAppBar(
                title = { Text(text = "Astral Agent") },
                actions = {
                    items[selected].actions(this)
                    actions()
                },
            )
        },
        content = { paddingValues ->
            NavHost(
                modifier = Modifier.padding(paddingValues),
                navController = navController,
                startDestination = items[selected].route,
            ) {
                items.forEach { item ->
                    composable(item.route) {
                        item.content()
                    }
                }
            }
        },
        bottomBar = {
            if (showBars) BottomNavigation {
                items.forEachIndexed { index, item ->
                    BottomNavigationItem(
                        selected = selected == index,
                        onClick = {
                            selected = index
                            onSelect(item)
                        },
                        icon = { Icon(imageVector = item.icon, contentDescription = "") },
                        label = { Text(text = item.label) }
                    )
                }
            }
        }
    )
}
