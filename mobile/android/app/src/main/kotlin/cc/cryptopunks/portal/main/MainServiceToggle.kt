package cc.cryptopunks.portal.main

//import androidx.compose.animation.core.Animatable
//import androidx.compose.animation.core.LinearEasing
//import androidx.compose.animation.core.RepeatMode
//import androidx.compose.animation.core.infiniteRepeatable
//import androidx.compose.animation.core.tween
//import androidx.compose.foundation.Image
//import androidx.compose.material.IconButton
//import androidx.compose.material.LocalContentColor
//import androidx.compose.runtime.Composable
//import androidx.compose.runtime.LaunchedEffect
//import androidx.compose.runtime.getValue
//import androidx.compose.runtime.mutableFloatStateOf
//import androidx.compose.runtime.mutableStateOf
//import androidx.compose.runtime.remember
//import androidx.compose.runtime.setValue
//import androidx.compose.ui.Modifier
//import androidx.compose.ui.draw.rotate
//import androidx.compose.ui.graphics.ColorFilter
//import androidx.compose.ui.platform.LocalContext
//import androidx.compose.ui.res.painterResource
//import androidx.compose.ui.tooling.preview.Preview
//import androidx.lifecycle.compose.collectAsStateWithLifecycle
//import cc.cryptopunks.portal.R
//import cc.cryptopunks.portal.node.AstralStatus
//import cc.cryptopunks.portal.node.astralStatus
//
//
//@Composable
//fun MainServiceToggle() {
//    val status by astralStatus.collectAsStateWithLifecycle()
//    val context = LocalContext.current
//    MainServiceToggle(status) {
//        when (status) {
//            AstralStatus.Stopped -> context.startAstralService()
//            AstralStatus.Started -> context.stopAstralService()
//            AstralStatus.Starting -> Unit
//        }
//    }
//}
//
//@Preview
//@Composable
//fun MainServiceTogglePreview() {
//    var status by remember { mutableStateOf(AstralStatus.Starting) }
//    MainServiceToggle(
//        status = status,
//        enabled = true,
//    ) {
//        val values = AstralStatus.values()
//        val next = (status.ordinal + 1) % values.size
//        status = values[next]
//    }
//}
//
//@Composable
//fun MainServiceToggle(
//    status: AstralStatus,
//    enabled: Boolean = status != AstralStatus.Starting,
//    onClick: () -> Unit,
//) {
//    IconButton(
//        onClick = onClick,
//        enabled = enabled,
//    ) {
//        val res = when (status) {
//            AstralStatus.Starting -> R.drawable.sync_black_24dp
//            AstralStatus.Started -> R.drawable.link_black_24dp
//            AstralStatus.Stopped -> R.drawable.link_off_black_24dp
//        }
//        var angle by remember(status) { mutableFloatStateOf(0f) }
//        val rotation = remember(status) { Animatable(angle) }
//        LaunchedEffect(status) {
//            if (status == AstralStatus.Starting) {
//                rotation.animateTo(
//                    targetValue = angle - 360f,
//                    animationSpec = infiniteRepeatable(
//                        animation = tween(920, easing = LinearEasing),
//                        repeatMode = RepeatMode.Restart
//                    )
//                ) {
//                    angle = value
//                }
//            }
//        }
//        Image(
//            painter = painterResource(res),
//            contentDescription = status.name,
//            modifier = Modifier.rotate(rotation.value),
//            colorFilter = ColorFilter.tint(LocalContentColor.current)
//        )
//    }
//}
