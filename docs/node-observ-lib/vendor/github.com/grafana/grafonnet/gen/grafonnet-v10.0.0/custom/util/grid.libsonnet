local d = import 'github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet';

{
  local root = self,

  local rowPanelHeight = 1,
  local gridWidth = 24,

  // Calculates the number of rows for a set of panels.
  countRows(panels, panelWidth):
    std.ceil(std.length(panels) / std.floor(gridWidth / panelWidth)),

  // Calculates gridPos for a panel based on its index, width and height.
  gridPosForIndex(index, panelWidth, panelHeight, startY): {
    local panelsPerRow = std.floor(gridWidth / panelWidth),
    local row = std.floor(index / panelsPerRow),
    local col = std.mod(index, panelsPerRow),
    gridPos: {
      w: panelWidth,
      h: panelHeight,
      x: panelWidth * col,
      y: startY + (panelHeight * row) + row,
    },
  },

  // Configures gridPos for each panel in a grid with equal width and equal height.
  makePanelGrid(panels, panelWidth, panelHeight, startY):
    std.mapWithIndex(
      function(i, panel)
        panel + root.gridPosForIndex(i, panelWidth, panelHeight, startY),
      panels
    ),

  '#makeGrid':: d.func.new(
    |||
      `makeGrid` returns an array of `panels` organized in a grid with equal `panelWidth`
      and `panelHeight`. Row panels are used as "linebreaks", if a Row panel is collapsed,
      then all panels below it will be folded into the row.

      This function will use the full grid of 24 columns, setting `panelWidth` to a value
      that can divide 24 into equal parts will fill up the page nicely. (1, 2, 3, 4, 6, 8, 12)
      Other value for `panelWidth` will leave a gap on the far right.

      Optional `startY` can be provided to place generated grid above or below existing panels.
    |||,
    args=[
      d.arg('panels', d.T.array),
      d.arg('panelWidth', d.T.number),
      d.arg('panelHeight', d.T.number),
      d.arg('startY', d.T.number),
    ],
  ),
  makeGrid(panels, panelWidth=8, panelHeight=8, startY=0):
    // Get indexes for all Row panels
    local rowIndexes = [
      i
      for i in std.range(0, std.length(panels) - 1)
      if panels[i].type == 'row'
    ];

    // Group panels below each Row panel
    local rowGroups =
      std.mapWithIndex(
        function(i, r) {
          header:
            {
              // Set initial values to ensure a value is set
              // may be overridden at per Row panel
              collapsed: false,
              panels: [],
            }
            + panels[r],
          panels:
            self.header.panels  // prepend panels that are part of the Row panel
            + (if i == std.length(rowIndexes) - 1  // last rowIndex
               then panels[r + 1:]
               else panels[r + 1:rowIndexes[i + 1]]),
          rows: root.countRows(self.panels, panelWidth),
        },
        rowIndexes
      );

    // Loop over rowGroups
    std.foldl(
      function(acc, rowGroup) acc + {
        local y = acc.nexty,
        nexty: y  // previous y
               + (rowGroup.rows * panelHeight)  // height of all rows
               + rowGroup.rows  // plus 1 for each row
               + acc.lastRowPanelHeight,

        lastRowPanelHeight: rowPanelHeight,  // set height for next round

        // Create a grid per group
        local panels = root.makePanelGrid(rowGroup.panels, panelWidth, panelHeight, y + 1),

        panels+:
          [
            // Add row header aka the Row panel
            rowGroup.header + {
              gridPos: {
                w: gridWidth,  // always full length
                h: rowPanelHeight,  // always 1 height
                x: 0,  // always at beginning
                y: y,
              },
              panels:
                // If row is collapsed, then store panels inside Row panel
                if rowGroup.header.collapsed
                then panels
                else [],
            },
          ]
          + (
            // If row is not collapsed, then expose panels directly
            if !rowGroup.header.collapsed
            then panels
            else []
          ),
      },
      rowGroups,
      {
        // Get panels that come before the rowGroups
        local panelsBeforeRowGroups =
          if std.length(rowIndexes) != 0
          then panels[0:rowIndexes[0]]
          else panels,  // matches all panels if no Row panels found
        local rows = root.countRows(panelsBeforeRowGroups, panelWidth),
        nexty: startY + (rows * panelHeight) + rows,

        lastRowPanelHeight: 0,  // starts without a row panel

        // Create a grid for the panels that come before the rowGroups
        panels: root.makePanelGrid(panelsBeforeRowGroups, panelWidth, panelHeight, startY),
      }
    ).panels,

  '#wrapPanels':: d.func.new(
    |||
      `wrapPanels` returns an array of `panels` organized in a grid, wrapping up to next 'row' if total width exceeds full grid of 24 columns.
      'panelHeight' and 'panelWidth' are used unless panels already have height and width defined.
    |||,
    args=[
      d.arg('panels', d.T.array),
      d.arg('panelWidth', d.T.number),
      d.arg('panelHeight', d.T.number),
      d.arg('startY', d.T.number),
    ],
  ),
  wrapPanels(panels, panelWidth=8, panelHeight=8, startY=0):
    std.foldl(
      function(acc, panel)
        if panel.type == 'row'
        then
          // when type=row, start new row immediatly and shift Y of new row by max height recorded
          acc + {
            panels+: [
              panel + {
                gridPos+:
                  {
                    x: acc.cursor.x,
                    y: acc.cursor.y + acc.cursor.maxH,
                    w: 0,
                    h: 1,
                  },
              },
            ],
            cursor:: {
              x: 0,
              y: acc.cursor.y + acc.cursor.maxH + 1,
              maxH: 0,
            },
          }
        else
          // handle regular panel
          local gridPos = std.get(panel, 'gridPos', {});
          local width = std.get(gridPos, 'w', panelWidth);
          local height = std.get(gridPos, 'h', panelHeight);
          if acc.cursor.x + width > gridWidth
          then
            // start new row as width exceeds gridWidth
            acc + {
              panels+: [
                panel + {
                  gridPos+:
                    {
                      x: 0,
                      y: acc.cursor.y + height,
                      w: width,
                      h: height,
                    },
                },
              ],
              cursor+:: {
                x: 0 + width,
                y: acc.cursor.y + height,
                maxH: if height > super.maxH then height else super.maxH,
              },
            }
          else
            // enough width, place panel on current row
            acc + {
              panels+: [
                panel + {
                  gridPos+:
                    {
                      x: acc.cursor.x,
                      y: acc.cursor.y,
                      w: width,
                      h: height,
                    },
                },
              ],
              cursor+:: {
                x: acc.cursor.x + width,
                y: acc.cursor.y,
                maxH: if height > super.maxH then height else super.maxH,
              },
            },
      panels,
      // Initial value for acc
      {
        panels: [],
        cursor:: {
          x: 0,
          y: startY,
          maxH: 0,  // max height of current 'row'
        },
      }
    ).panels,
}
